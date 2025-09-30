package utils

import "github.com/niko-2609/tracker-expense/database"

func UpdateDashboardMetrics(userID uint) error {
	// A. Total Income / Expense / Net Savings
	totalQuery := `
    WITH totals AS (
        SELECT
            SUM(CASE WHEN txn_type='income' THEN amount ELSE 0 END) AS total_income,
            SUM(CASE WHEN txn_type='expense' THEN amount ELSE 0 END) AS total_expense
        FROM transactions
        WHERE user_id = ?
    )
    INSERT INTO user_dashboard_metrics (user_id, total_income, total_expense, net_savings, updated_at)
    SELECT ?, total_income, total_expense, total_income - total_expense, NOW()
    FROM totals
    ON CONFLICT (user_id) DO UPDATE
    SET total_income = EXCLUDED.total_income,
        total_expense = EXCLUDED.total_expense,
        net_savings = EXCLUDED.net_savings,
        updated_at = NOW();
    `
	if err := database.DB.Exec(totalQuery, userID, userID).Error; err != nil {
		return err
	}

	// B. Monthly Totals (for line chart)
	monthlyQuery := `
	WITH monthly AS (
	    SELECT
	        TO_CHAR(DATE_TRUNC('month', txn_date), 'YYYY-MM') AS month,
	        SUM(CASE WHEN txn_type='income' THEN amount ELSE -amount END) AS net
	    FROM transactions
	    WHERE user_id = ?
	    GROUP BY DATE_TRUNC('month', txn_date)
	    ORDER BY month
	)
	UPDATE user_dashboard_metrics
	SET monthly_totals = (SELECT jsonb_object_agg(month, net) FROM monthly)
	WHERE user_id = ?;
	`
	if err := database.DB.Exec(monthlyQuery, userID, userID).Error; err != nil {
		return err
	}

	// // C. Top 5 Expense Categories (for pie chart)
	topCatQuery := `
	UPDATE user_dashboard_metrics
	SET top_expense_categories = (
    	SELECT jsonb_agg(jsonb_build_object('category', category, 'amount', amount))
    		FROM (
        	SELECT c.name AS category, SUM(t.amount) AS amount
        	FROM transactions t
        	JOIN categories c ON t.category_id = c.id
        	WHERE t.user_id = ? AND t.txn_type = 'expense'
        	GROUP BY c.id, c.name
        	ORDER BY SUM(t.amount) DESC
        	LIMIT 5
    	) sub
	)
	WHERE user_id = ?;
	`
	if err := database.DB.Exec(topCatQuery, userID, userID).Error; err != nil {
		return err
	}

	return nil
}
