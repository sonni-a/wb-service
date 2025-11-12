package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	emailRegex = regexp.MustCompile(`^[\w._%+\-]+@[\w.\-]+\.[A-Za-z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[0-9][0-9\s\-\(\)]{6,19}$`)
	zipRegex   = regexp.MustCompile(`^\d{4,10}$`)
)

func (d *Delivery) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("delivery name cannot be empty")
	}
	if strings.TrimSpace(d.Phone) == "" {
		return errors.New("delivery phone cannot be empty")
	}
	if strings.TrimSpace(d.Zip) == "" {
		return errors.New("delivery zip cannot be empty")
	}
	if strings.TrimSpace(d.City) == "" {
		return errors.New("delivery city cannot be empty")
	}
	if strings.TrimSpace(d.Address) == "" {
		return errors.New("delivery address cannot be empty")
	}
	if strings.TrimSpace(d.Region) == "" {
		return errors.New("delivery region cannot be empty")
	}
	if strings.TrimSpace(d.Email) == "" {
		return errors.New("delivery email cannot be empty")
	}

	if !phoneRegex.MatchString(d.Phone) {
		return fmt.Errorf("invalid phone format: %s", d.Phone)
	}
	if !emailRegex.MatchString(d.Email) {
		return fmt.Errorf("invalid email format: %s", d.Email)
	}
	if !zipRegex.MatchString(d.Zip) {
		return fmt.Errorf("invalid zip code: %s", d.Zip)
	}

	return nil
}

var (
	currencyRegex = regexp.MustCompile(`^[A-Z]{3}$`)
	providerRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{2,}$`)
	bankRegex     = regexp.MustCompile(`^[a-zA-Z0-9\s-]{2,50}$`)
)

func (p *Payment) Validate() error {
	if p == nil {
		return errors.New("payment cannot be nil")
	}

	if strings.TrimSpace(p.Transaction) == "" {
		return errors.New("transaction cannot be empty")
	}

	if !currencyRegex.MatchString(p.Currency) {
		return fmt.Errorf("invalid currency: %s", p.Currency)
	}

	if strings.TrimSpace(p.Provider) == "" {
		return errors.New("provider cannot be empty")
	}
	if !providerRegex.MatchString(p.Provider) {
		return fmt.Errorf("invalid provider format: %s", p.Provider)
	}

	if p.Amount <= 0 {
		return fmt.Errorf("invalid payment amount: %d", p.Amount)
	}
	if p.DeliveryCost < 0 || p.GoodsTotal < 0 || p.CustomFee < 0 {
		return errors.New("cost fields must be non-negative")
	}

	if p.PaymentDt <= 0 {
		return errors.New("invalid payment timestamp")
	}
	if time.Unix(p.PaymentDt, 0).After(time.Now().Add(24 * time.Hour)) {
		return errors.New("payment date cannot be in the future")
	}

	if strings.TrimSpace(p.Bank) == "" {
		return errors.New("bank cannot be empty")
	}
	if !bankRegex.MatchString(p.Bank) {
		return fmt.Errorf("invalid bank name: %s", p.Bank)
	}

	return nil
}

func (i *Item) Validate() error {
	if strings.TrimSpace(i.OrderUID) == "" {
		return errors.New("item order_uid cannot be empty")
	}
	if i.ChrtID == 0 {
		return errors.New("item chrt_id cannot be zero")
	}
	if strings.TrimSpace(i.TrackNumber) == "" {
		return errors.New("item track_number cannot be empty")
	}
	if i.Price <= 0 {
		return errors.New("item price must be positive")
	}
	if strings.TrimSpace(i.RID) == "" {
		return errors.New("item rid cannot be empty")
	}
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if i.Sale < 0 {
		return errors.New("item sale cannot be negative")
	}
	if strings.TrimSpace(i.Size) == "" {
		return errors.New("item size cannot be empty")
	}
	if i.TotalPrice <= 0 {
		return errors.New("item total_price must be positive")
	}
	if i.NmID == 0 {
		return errors.New("item nm_id cannot be zero")
	}
	if strings.TrimSpace(i.Brand) == "" {
		return errors.New("item brand cannot be empty")
	}
	if i.Status <= 0 {
		return fmt.Errorf("item status must be positive, got %d", i.Status)
	}
	return nil
}

func (o *Order) Validate() error {
	if strings.TrimSpace(o.OrderUID) == "" {
		return errors.New("order_uid cannot be empty")
	}
	if strings.TrimSpace(o.TrackNumber) == "" {
		return errors.New("track_number cannot be empty")
	}
	if strings.TrimSpace(o.Entry) == "" {
		return errors.New("entry cannot be empty")
	}
	if strings.TrimSpace(o.Locale) == "" {
		return errors.New("locale cannot be empty")
	}
	if strings.TrimSpace(o.CustomerID) == "" {
		return errors.New("customer_id cannot be empty")
	}
	if strings.TrimSpace(o.DeliveryService) == "" {
		return errors.New("delivery_service cannot be empty")
	}
	if strings.TrimSpace(o.ShardKey) == "" {
		return errors.New("shardkey cannot be empty")
	}
	if o.SmID <= 0 {
		return errors.New("sm_id must be positive")
	}
	if strings.TrimSpace(o.OofShard) == "" {
		return errors.New("oof_shard cannot be empty")
	}

	if err := o.Delivery.Validate(); err != nil {
		return err
	}

	if err := o.Payment.Validate(); err != nil {
		return err
	}

	if len(o.Items) == 0 {
		return errors.New("items cannot be empty")
	}
	for i, item := range o.Items {
		if err := item.Validate(); err != nil {
			return errors.New("item[" + strconv.Itoa(i) + "]: " + err.Error())
		}
	}

	return nil
}
