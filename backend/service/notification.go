package service

import (
	"context"
	"log"
	"time"
	"todo/backend/db"

	"github.com/gen2brain/beeep"
)

// StartNotificationScheduler kicks off the periodic reminder check. The
// caller owns ctx and must cancel it (typically in the Wails OnShutdown hook)
// so the goroutine exits cleanly when the app quits.
func StartNotificationScheduler(ctx context.Context) {
	// Check immediately on start
	go checkReminders()

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				checkReminders()
			}
		}
	}()
}

// sendNotification is the actual desktop-notification dispatch. It's a package
// variable so tests can stub it without invoking beeep (which blocks on dbus
// in headless environments).
var sendNotification = func(title, description string) error {
	log.Printf("Sending notification for task: %s", title)
	return beeep.Notify("Todo Reminder", title, "")
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

	// Only pick tasks that haven't been notified yet. Mark as notified right after
	// dispatching so a stuck window won't re-deliver the same reminder every minute.
	rows, err := db.DB.Query("SELECT id, title, description FROM todos WHERE completed = false AND remind_at >= ? AND remind_at < ? AND notified_at IS NULL AND deleted_at IS NULL", start, end)
	if err != nil {
		log.Println("Error checking reminders:", err)
		return
	}

	// Drain rows into a slice and close immediately. Holding rows open while
	// we run UPDATE / beeep.Notify would deadlock the sql.DB connection pool:
	// the outer Query is still using a connection, and any Exec inside the
	// loop needs another one.
	type dueReminder struct {
		id                       int
		title, description       string
	}
	var due []dueReminder
	for rows.Next() {
		var r dueReminder
		if err := rows.Scan(&r.id, &r.title, &r.description); err != nil {
			continue
		}
		due = append(due, r)
	}
	rows.Close()

	for _, r := range due {
		if err := sendNotification(r.title, r.description); err != nil {
			log.Println("Error sending notification:", err)
			continue
		}

		// Mark as notified to prevent re-sending while the window slides forward.
		if _, err := db.DB.Exec("UPDATE todos SET notified_at = ? WHERE id = ?", time.Now().UTC(), r.id); err != nil {
			log.Println("Error marking todo as notified:", err)
		}
	}
}
