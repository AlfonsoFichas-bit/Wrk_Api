package handlers

import (
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"Wrk_Api/internal/database"
	"Wrk_Api/internal/models"

	"github.com/gin-gonic/gin"
)

// Burndown Logic
func GetSprintBurndown(c *gin.Context) {
	sprintID := c.Param("sprintId")
	var sprint models.Sprint
	if err := database.DB.Preload("UserStories").First(&sprint, "id = ?", sprintID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sprint not found"})
		return
	}

	totalPoints := 0
	for _, story := range sprint.UserStories {
		if story.StoryPoints != nil {
			totalPoints += *story.StoryPoints
		}
	}

	// Simple Ideal vs Actual calculation logic
	// In a real app, we need to track *when* each story was completed daily.
	// We can use CompletedAt field.
	
	type DataPoint struct {
		Day    int     `json:"day"`
		Date   string  `json:"date"`
		Ideal  float64 `json:"ideal"`
		Actual *int    `json:"actual"`
	}
	
	series := []DataPoint{}
	
	if !sprint.StartDate.IsZero() && !sprint.EndDate.IsZero() {
		days := int(sprint.EndDate.Sub(sprint.StartDate).Hours() / 24)
		idealDec := float64(totalPoints) / float64(days)
		
		for i := 0; i <= days; i++ {
			date := sprint.StartDate.Add(time.Hour * 24 * time.Duration(i))
			ideal := math.Max(0, float64(totalPoints)-(idealDec*float64(i)))
			
			// Actual calculation: Total - (Sum of points completed <= date)
			burned := 0
			for _, story := range sprint.UserStories {
				if story.CompletedAt != nil && !story.CompletedAt.After(date) {
					if story.StoryPoints != nil {
						burned += *story.StoryPoints
					}
				}
			}
			
			var actual *int
			if date.Before(time.Now().Add(time.Hour * 24)) {
				rem := totalPoints - burned
				actual = &rem
			}
			
			series = append(series, DataPoint{
				Day:    i,
				Date:   date.Format("2006-01-02"),
				Ideal:  ideal,
				Actual: actual,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"totalPoints": totalPoints, "series": series}})
}

// Velocity Logic
func GetProjectVelocity(c *gin.Context) {
	projectID := c.Param("projectId")
	var sprints []models.Sprint
	if err := database.DB.Preload("UserStories").Where("project_id = ?", projectID).Order("start_date asc").Find(&sprints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sprints"})
		return
	}

	type VelocityData struct {
		Name      string `json:"name"`
		Committed int    `json:"committed"`
		Completed int    `json:"completed"`
	}
	
	data := []VelocityData{}
	
	for _, s := range sprints {
		committed := 0
		completed := 0
		for _, us := range s.UserStories {
			pts := 0
			if us.StoryPoints != nil {
				pts = *us.StoryPoints
			}
			committed += pts
			if us.CompletedAt != nil {
				completed += pts
			}
		}
		data = append(data, VelocityData{Name: s.Name, Committed: committed, Completed: completed})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Contribution Logic
func GetProjectContribution(c *gin.Context) {
	projectID := c.Param("projectId")
	
	type Contrib struct {
		User  models.User `json:"user"`
		Count int         `json:"count"`
	}
	
	// Count completed tasks by assignee
	rows, err := database.DB.Table("tasks").
		Select("assignee_id, count(*) as count").
		Where("project_id = ? AND status IN ?", projectID, []string{"COMPLETED", "DONE"}).
		Group("assignee_id").
		Rows()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating contribution"})
		return
	}
	defer rows.Close()

	results := []Contrib{}
	for rows.Next() {
		var assigneeID string
		var count int
		rows.Scan(&assigneeID, &count)
		
		if assigneeID != "" {
			var u models.User
			database.DB.First(&u, "id = ?", assigneeID)
			results = append(results, Contrib{User: u, Count: count})
		}
	}
	
	// Sort desc
	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func ExportProjectCSV(c *gin.Context) {
	projectID := c.Param("projectId")
	
	var project models.Project
	if err := database.DB.Preload("Sprints.Tasks.Assignee").First(&project, "id = ?", projectID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	csv := "Sprint,Task,Assignee,Status,Priority\n"
	for _, sprint := range project.Sprints {
		for _, task := range sprint.Tasks {
			assignee := "Unassigned"
			if task.Assignee != nil {
				assignee = task.Assignee.Name
			}
			csv += fmt.Sprintf("%s,%s,%s,%s,%s\n", sprint.Name, task.Title, assignee, task.Status, task.Priority)
		}
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"project-%s.csv\"", projectID))
	c.String(http.StatusOK, csv)
}
