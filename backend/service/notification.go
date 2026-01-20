package service

import (
	"log"
	"time"
	"todo/backend/db"

	"github.com/gen2brain/beeep"
)

func StartNotificationScheduler() {
	// Check immediately on start
	go checkReminders()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			checkReminders()
		}
	}()
}

func checkReminders() {
	// Use UTC because backend stores times in UTC (from ISO strings)
	now := time.Now().UTC()
	// Check for reminders in the current minute window (or slightly past if we missed it)
	// Let's look for anything scheduled in the [now, now+1m) window.
	// But to be safe against slight skews, maybe [now-30s, now+30s]?
	// Actually, best is [now, now+1m) for forward looking.
	
	start := now
	end := now.Add(1 * time.Minute)

	rows, err := db.DB.Query("SELECT title, description FROM todos WHERE completed = false AND remind_at >= ? AND remind_at < ?", start, end)
	if err != nil {
		log.Println("Error checking reminders:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var title, description string
		if err := rows.Scan(&title, &description); err != nil {
			continue
		}
		
		// Send notification
		log.Printf("Sending notification for task: %s", title)
		err := beeep.Notify("Todo Reminder", title, "")
		if err != nil {
			log.Println("Error sending notification:", err)
		}
	}
}
