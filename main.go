package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// CalendarWidget represents a custom calendar widget
type CalendarWidget struct {
	widget.BaseWidget
	currentTime    time.Time
	selectedTime   time.Time
	onDateSelected func(time.Time)
}

func NewCalendarWidget(initialTime time.Time, onDateSelected func(time.Time)) *CalendarWidget {
	c := &CalendarWidget{
		currentTime:    initialTime,
		selectedTime:   initialTime,
		onDateSelected: onDateSelected,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *CalendarWidget) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	return &calendarRenderer{
		calendar: c,
		objects:  []fyne.CanvasObject{},
	}
}

func (c *CalendarWidget) SetMonth(t time.Time) {
	c.currentTime = t
	c.Refresh()
}

func (c *CalendarWidget) PreviousMonth() {
	c.currentTime = c.currentTime.AddDate(0, -1, 0)
	c.Refresh()
}

func (c *CalendarWidget) NextMonth() {
	c.currentTime = c.currentTime.AddDate(0, 1, 0)
	c.Refresh()
}

type calendarRenderer struct {
	calendar *CalendarWidget
	objects  []fyne.CanvasObject
}

func (r *calendarRenderer) Layout(size fyne.Size) {
	// Layout is handled by the container
}

func (r *calendarRenderer) MinSize() fyne.Size {
	return fyne.NewSize(350, 300)
}

func (r *calendarRenderer) Refresh() {
	r.calendar.BaseWidget.Refresh()
}

func (r *calendarRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *calendarRenderer) Destroy() {}

// CreateTimePicker creates a shadcn-style time picker
func CreateTimePicker(currentTime time.Time, onTimeSelected func(time.Time)) fyne.CanvasObject {
	// Use pointers to allow modification in closures
	hours := currentTime.Hour()
	minutes := currentTime.Minute()
	seconds := currentTime.Second()

	// Define updateTime before the select widgets
	updateTime := func() {
		selectedTime := time.Date(
			currentTime.Year(),
			currentTime.Month(),
			currentTime.Day(),
			hours,
			minutes,
			seconds,
			0,
			time.Local,
		)
		onTimeSelected(selectedTime)
	}

	// Create select widgets
	hourSelect := widget.NewSelect(
		generateRange(0, 23),
		func(value string) {
			if value != "" {
				var hour int
				fmt.Sscanf(value, "%d", &hour)
				hours = hour
				updateTime()
			}
		},
	)
	hourSelect.SetSelected(fmt.Sprintf("%02d", hours))

	minuteSelect := widget.NewSelect(
		generateRange(0, 59),
		func(value string) {
			if value != "" {
				var minute int
				fmt.Sscanf(value, "%d", &minute)
				minutes = minute
				updateTime()
			}
		},
	)
	minuteSelect.SetSelected(fmt.Sprintf("%02d", minutes))

	secondSelect := widget.NewSelect(
		generateRange(0, 59),
		func(value string) {
			if value != "" {
				var second int
				fmt.Sscanf(value, "%d", &second)
				seconds = second
				updateTime()
			}
		},
	)
	secondSelect.SetSelected(fmt.Sprintf("%02d", seconds))

	// Time display with colons
	timeContainer := container.NewGridWithColumns(5,
		hourSelect,
		widget.NewLabel(":"),
		minuteSelect,
		widget.NewLabel(":"),
		secondSelect,
	)

	return container.NewVBox(
		widget.NewLabel("Select Time:"),
		timeContainer,
	)
}

func generateRange(min, max int) []string {
	result := make([]string, max-min+1)
	for i := range result {
		result[i] = fmt.Sprintf("%02d", min+i)
	}
	return result
}

// CreateCalendarView creates a calendar view widget
func CreateCalendarView(initialTime time.Time, onDateSelected func(time.Time)) fyne.CanvasObject {
	currentMonth := initialTime

	// Month/Year label
	monthYearLabel := widget.NewLabel("")
	monthYearLabel.Alignment = fyne.TextAlignCenter
	monthYearLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Navigation buttons
	prevBtn := widget.NewButton("<", nil)
	nextBtn := widget.NewButton(">", nil)

	// Calendar grid
	calendarGrid := container.NewGridWithColumns(7)

	updateCalendar := func() {
		// Update month/year label
		monthYearLabel.SetText(currentMonth.Format("January 2006"))

		// Clear grid
		calendarGrid.Objects = nil

		// Add day headers
		days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
		for _, day := range days {
			dayLabel := widget.NewLabel(day)
			dayLabel.Alignment = fyne.TextAlignCenter
			calendarGrid.Add(dayLabel)
		}

		// Get first day of month and number of days
		firstOfMonth := time.Date(currentMonth.Year(), currentMonth.Month(), 1, 0, 0, 0, 0, time.Local)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		firstDayWeekday := int(firstOfMonth.Weekday())
		daysInMonth := lastOfMonth.Day()

		// Add empty cells for days before the first of the month
		for i := 0; i < firstDayWeekday; i++ {
			calendarGrid.Add(widget.NewLabel(""))
		}

		// Add day buttons
		today := time.Now()
		for day := 1; day <= daysInMonth; day++ {
			dayTime := time.Date(currentMonth.Year(), currentMonth.Month(), day, 0, 0, 0, 0, time.Local)
			isToday := dayTime.Year() == today.Year() &&
				dayTime.Month() == today.Month() &&
				dayTime.Day() == today.Day()

			dayBtn := widget.NewButton(fmt.Sprintf("%d", day), nil)

			if isToday {
				dayBtn.Importance = widget.HighImportance
			}

			dayBtn.OnTapped = func() {
				onDateSelected(dayTime)
			}

			calendarGrid.Add(dayBtn)
		}

		calendarGrid.Refresh()
	}

	// Setup navigation
	prevBtn.OnTapped = func() {
		currentMonth = currentMonth.AddDate(0, -1, 0)
		updateCalendar()
	}

	nextBtn.OnTapped = func() {
		currentMonth = currentMonth.AddDate(0, 1, 0)
		updateCalendar()
	}

	// Initial calendar update
	updateCalendar()

	// Navigation header
	navHeader := container.NewBorder(nil, nil, prevBtn, nextBtn, monthYearLabel)

	// Combine everything
	calendarContainer := container.NewVBox(
		navHeader,
		widget.NewSeparator(),
		calendarGrid,
	)

	return calendarContainer
}

func main() {
	a := app.New()
	w := a.NewWindow("Time Changer")
	w.Resize(fyne.NewSize(500, 600))

	// Get current time
	currentTime := time.Now()

	// Selected date label
	selectedDateLabel := widget.NewLabel("Selected Date:")
	selectedDateLabel.TextStyle = fyne.TextStyle{Bold: true}

	selectedDateValue := widget.NewLabel(currentTime.Format("2006-01-02"))

	// Display label for selected datetime
	selectedDateTimeLabel := widget.NewLabel(fmt.Sprintf(
		"%04d-%02d-%02d %02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second(),
	))
	selectedDateTimeLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Update preview function (must be defined before it's used)
	updatePreview := func() {
		date, _ := time.Parse("2006-01-02", selectedDateValue.Text)
		selectedDateTimeLabel.SetText(fmt.Sprintf(
			"%04d-%02d-%02d %02d:%02d:%02d",
			date.Year(), date.Month(), date.Day(),
			currentTime.Hour(), currentTime.Minute(), currentTime.Second(),
		))
	}

	// Calendar with date selection callback
	calendar := CreateCalendarView(currentTime, func(selectedDate time.Time) {
		selectedDateValue.SetText(selectedDate.Format("2006-01-02"))
		currentTime = selectedDate
		updatePreview()
	})

	// Time picker with callback
	timePicker := CreateTimePicker(currentTime, func(selectedTime time.Time) {
		currentTime = selectedTime
		updatePreview()
	})

	// Set current time button
	setCurrentButton := widget.NewButton("Set to Current Time", func() {
		now := time.Now()
		currentTime = now
		selectedDateValue.SetText(now.Format("2006-01-02"))
		updatePreview()
	})

	// Layout
	content := container.NewVBox(
		selectedDateLabel,
		calendar,
		selectedDateValue,
		widget.NewSeparator(),
		timePicker,
		widget.NewSeparator(),
		setCurrentButton,
		widget.NewSeparator(),
		widget.NewLabel("Preview:"),
		selectedDateTimeLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
