package repositories

import (
	"bank-api/models"
	"database/sql"
	"errors"
	"time"
)

type CreditRepository struct {
	db *sql.DB
}

func NewCreditRepository(db *sql.DB) *CreditRepository {
	return &CreditRepository{db: db}
}

func (r *CreditRepository) Create(credit *models.Credit) error {
	query := `
		INSERT INTO credits (user_id, amount, rate, term_months, monthly_payment, total_payment, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		credit.UserID,
		credit.Amount,
		credit.Rate,
		credit.TermMonths,
		credit.MonthlyPayment,
		credit.TotalPayment,
		credit.Status,
	).Scan(&credit.ID, &credit.CreatedAt)
}

func (r *CreditRepository) FindByUserID(userID int) ([]models.Credit, error) {
	query := `SELECT id, user_id, amount, rate, term_months, monthly_payment, total_payment, status, created_at FROM credits WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var credits []models.Credit
	for rows.Next() {
		var c models.Credit
		if err := rows.Scan(&c.ID, &c.UserID, &c.Amount, &c.Rate, &c.TermMonths,
			&c.MonthlyPayment, &c.TotalPayment, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}
	return credits, rows.Err()
}

func (r *CreditRepository) FindByID(id int) (*models.Credit, error) {
	credit := &models.Credit{}
	query := `SELECT id, user_id, amount, rate, term_months, monthly_payment, total_payment, status, created_at FROM credits WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&credit.ID, &credit.UserID, &credit.Amount, &credit.Rate,
		&credit.TermMonths, &credit.MonthlyPayment, &credit.TotalPayment,
		&credit.Status, &credit.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("credit not found")
		}
		return nil, err
	}
	return credit, nil
}

func (r *CreditRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE credits SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

// Payment Schedule methods
func (r *CreditRepository) CreatePaymentSchedule(schedule *models.PaymentSchedule) error {
	query := `
		INSERT INTO payment_schedules (credit_id, due_date, amount, paid, penalty)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return r.db.QueryRow(
		query,
		schedule.CreditID,
		schedule.DueDate,
		schedule.Amount,
		schedule.Paid,
		schedule.Penalty,
	).Scan(&schedule.ID)
}

func (r *CreditRepository) FindScheduleByCreditID(creditID int) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, due_date, amount, paid, paid_at, penalty FROM payment_schedules WHERE credit_id = $1 ORDER BY due_date`

	rows, err := r.db.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.DueDate, &s.Amount, &s.Paid, &s.PaidAt, &s.Penalty); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

func (r *CreditRepository) FindOverduePayments() ([]models.PaymentSchedule, error) {
	query := `
		SELECT id, credit_id, due_date, amount, paid, paid_at, penalty
		FROM payment_schedules
		WHERE paid = FALSE AND due_date < $1
		ORDER BY due_date`

	rows, err := r.db.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.DueDate, &s.Amount, &s.Paid, &s.PaidAt, &s.Penalty); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, rows.Err()
}

func (r *CreditRepository) MarkPaymentPaid(scheduleID int, paidAt time.Time) error {
	query := `UPDATE payment_schedules SET paid = TRUE, paid_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, paidAt, scheduleID)
	return err
}

func (r *CreditRepository) AddPenalty(scheduleID int, penalty float64) error {
	query := `UPDATE payment_schedules SET penalty = penalty + $1 WHERE id = $2`
	_, err := r.db.Exec(query, penalty, scheduleID)
	return err
}