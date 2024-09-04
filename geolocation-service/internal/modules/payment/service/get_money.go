package service

import (
	"context"
	"geolocation-service/internal/models"
	"time"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/sub"
)

type PaymentService struct {
	apiKey    string
	productID string
	//paymentMethodID string
	trialEnd int64
	stor     PayStorager
}

func NewPaymentService(apiKey, productID string, trialEnd int64, stor PayStorager) *PaymentService {
	stripe.Key = apiKey
	return &PaymentService{
		apiKey:    apiKey,
		productID: productID,
		//	paymentMethodID: paymentMethodID,
		trialEnd: trialEnd,
		stor:     stor,
	}
}

type PayStorager interface {
	AddCustomer(ctx context.Context, user *models.User) error
	UpdateCustomer(ctx context.Context, user *models.User) error
	GetCustomer(ctx context.Context, userID string) (*models.User, error)
	DeleteCustomer(ctx context.Context, userID string) error
}

// Stripe предоставляет несколько статусов для подписок, которые помогают отслеживать состояние подписки клиента.
// Вот основные статусы подписки, которые вы можете получить с помощью Stripe API:
// incomplete: Подписка создана, но первый платеж не был успешно завершен. Это может произойти, если платеж требует дополнительных шагов, таких как 3-D Secure.
// incomplete_expired: Подписка была в статусе incomplete более 23 часов, и клиент не был выставлен счет. Подписка автоматически отменяется.
// trialing: Подписка находится в пробном периоде. Этот статус будет действовать до окончания пробного периода.
// active: Подписка активна и платежи успешно обрабатываются.
// past_due: Последний счет не был оплачен вовремя. Stripe будет пытаться повторно списать платеж в соответствии с вашими настройками.
// canceled: Подписка была отменена и больше не активна.
// unpaid: Подписка не оплачена, и счета остаются открытыми, но платежи не будут пытаться списываться до добавления нового метода оплаты.

func (p *PaymentService) CreateCustomer(user *models.User) (*models.User, error) {
	ctx := context.Background()
	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.Name),
	}
	c, err := customer.New(params)
	if err != nil {
		return nil, err
	}
	user.CustomerID = c.ID
	err = p.stor.AddCustomer(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func createPlanWithProductID(productID string) (string, error) {
	stripe.Key = "your-stripe-secret-key"

	params := &stripe.PlanParams{
		Amount:   stripe.Int64(5000), // Сумма в центах (например, $50.00)
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Interval: stripe.String("month"), // Интервал подписки (например, ежемесячно)
		Product: &stripe.PlanProductParams{
			ID: stripe.String(productID),
		},
	}

	p, err := plan.New(params)
	if err != nil {
		return "", err
	}

	return p.ID, nil
}

func (p *PaymentService) CreateSubscription(user *models.User) error {
	// Создаем план с использованием Product ID
	planID, err := createPlanWithProductID(p.productID)
	if err != nil {
		return err
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(user.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(planID),
			},
		},
		TrialEnd: stripe.Int64(p.trialEnd),
		// TODO DefaultPaymentMethod: stripe.String(paymentMethodID),
	}
	_, err = sub.New(subscriptionParams)
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentService) GetSubscriptionEndDate(user *models.User) (time.Time, error) {
	ctx := context.Background()
	s, err := sub.Get(user.CustomerID, nil)
	if err != nil {
		return time.Time{}, err
	}
	if s.CurrentPeriodEnd == 0 {
		return time.Time{}, nil
	}

	user.SubscriptionEndDate = time.Unix(s.CurrentPeriodEnd, 0)

	err = p.stor.UpdateCustomer(ctx, user)
	if err != nil {
		return time.Time{}, err
	}
	return user.SubscriptionEndDate, nil
}

func (p *PaymentService) GetSubscriptionStatus(user *models.User) (string, error) {
	ctx := context.Background()
	s, err := sub.Get(user.CustomerID, nil)
	if err != nil {
		return "", err
	}

	user.StatusSubscription = string(s.Status)
	err = p.stor.UpdateCustomer(ctx, user)
	if err != nil {
		return "", err
	}

	return string(s.Status), nil
}
