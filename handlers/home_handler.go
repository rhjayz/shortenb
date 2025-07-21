package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"shortb/config"
	"shortb/middleware"
	"shortb/models"
	"shortb/utils"
)

type MonthlyClick struct {
	Month string
	Total int
}

type HomeData struct {
	LabelsJSON template.JS
	ChartJSON  template.JS
}

func Index(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r)
	if !ok {
		user = models.User{}
	}

	var createdLink, customDomain, totalClick int

	row := config.DB.Raw(`
	SELECT 
		COUNT(DISTINCT shortlinks.short_code),
		COUNT(DISTINCT CASE 
			WHEN shortlinks.custom_domain IS NOT NULL AND shortlinks.custom_domain != '' 
			THEN shortlinks.id 
		END),
		COUNT(clicklogs.id)
	FROM clicklogs
	JOIN shortlinks ON shortlinks.id = clicklogs.shortlink_id
	WHERE shortlinks.user_id = ?
	  AND clicklogs.deleted_at IS NULL
	`, user.ID).Row()

	err := row.Scan(&createdLink, &customDomain, &totalClick)
	if err != nil {
		fmt.Println("Scan error:", err)
	}

	var clicks []MonthlyClick

	err = config.DB.
		Raw(`
		SELECT TO_CHAR(clicked_at, 'Mon') AS month, COUNT(*) AS total
		FROM clicklogs
		JOIN shortlinks ON shortlinks.id = clicklogs.shortlink_id
		WHERE shortlinks.user_id = ? AND clicked_at >= date_trunc('month', now() - interval '1 month')
		GROUP BY TO_CHAR(clicked_at, 'Mon'), date_trunc('month', clicked_at)
		ORDER BY date_trunc('month', clicked_at)
	`, user.ID).
		Scan(&clicks).Error

	if err != nil {
		fmt.Println("Scan error:", err)
	}

	labels := []string{}
	chart := []int{}

	for _, c := range clicks {
		labels = append(labels, c.Month)
		chart = append(chart, c.Total)
	}

	labelJSON, _ := json.Marshal(labels)
	chartJSON, _ := json.Marshal(chart)

	data := map[string]interface{}{
		"User":       user,
		"Created":    createdLink,
		"Click":      totalClick,
		"Custom":     customDomain,
		"LabelsJSON": template.JS(labelJSON),
		"ChartJSON":  template.JS(chartJSON),
	}

	utils.RenderTemplate(w, "views/pages/home/index.html", data)
}
