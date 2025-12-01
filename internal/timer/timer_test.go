package timer

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	d := 5 * time.Minute
	timer := New(d)

	if timer.Duration != d {
		t.Errorf("Duration = %v, want %v", timer.Duration, d)
	}
	if timer.Remaining != d {
		t.Errorf("Remaining = %v, want %v", timer.Remaining, d)
	}
	if timer.State != Stopped {
		t.Errorf("State = %v, want %v", timer.State, Stopped)
	}
}

func TestStart(t *testing.T) {
	tests := []struct {
		name         string
		initialState State
	}{
		{"from Stopped", Stopped},
		{"from Paused", Paused},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(time.Minute)
			timer.State = tt.initialState
			timer.Start()
			if timer.State != Running {
				t.Errorf("State = %v, want %v", timer.State, Running)
			}
		})
	}
}

func TestPause(t *testing.T) {
	tests := []struct {
		name          string
		initialState  State
		expectedState State
	}{
		{"from Running", Running, Paused},
		{"from Stopped", Stopped, Stopped},
		{"from Paused", Paused, Paused},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(time.Minute)
			timer.State = tt.initialState
			timer.Pause()
			if timer.State != tt.expectedState {
				t.Errorf("State = %v, want %v", timer.State, tt.expectedState)
			}
		})
	}
}

func TestToggle(t *testing.T) {
	tests := []struct {
		name          string
		initialState  State
		expectedState State
	}{
		{"Running to Paused", Running, Paused},
		{"Paused to Running", Paused, Running},
		{"Stopped to Running", Stopped, Running},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(time.Minute)
			timer.State = tt.initialState
			timer.Toggle()
			if timer.State != tt.expectedState {
				t.Errorf("State = %v, want %v", timer.State, tt.expectedState)
			}
		})
	}
}

func TestReset(t *testing.T) {
	timer := New(5 * time.Minute)
	timer.Remaining = time.Minute
	timer.State = Running

	timer.Reset()

	if timer.Remaining != timer.Duration {
		t.Errorf("Remaining = %v, want %v", timer.Remaining, timer.Duration)
	}
	if timer.State != Stopped {
		t.Errorf("State = %v, want %v", timer.State, Stopped)
	}
}

func TestResetWith(t *testing.T) {
	timer := New(5 * time.Minute)
	timer.Remaining = time.Minute
	timer.State = Running

	newDuration := 10 * time.Minute
	timer.ResetWith(newDuration)

	if timer.Duration != newDuration {
		t.Errorf("Duration = %v, want %v", timer.Duration, newDuration)
	}
	if timer.Remaining != newDuration {
		t.Errorf("Remaining = %v, want %v", timer.Remaining, newDuration)
	}
	if timer.State != Stopped {
		t.Errorf("State = %v, want %v", timer.State, Stopped)
	}
}

func TestTick(t *testing.T) {
	tests := []struct {
		name              string
		initialState      State
		initialRemaining  time.Duration
		tickAmount        time.Duration
		expectedRemaining time.Duration
		expectedComplete  bool
		expectedState     State
	}{
		{"running normal tick", Running, 5 * time.Second, time.Second, 4 * time.Second, false, Running},
		{"running completes", Running, time.Second, time.Second, 0, true, Stopped},
		{"running overshoots", Running, time.Second, 2 * time.Second, 0, true, Stopped},
		{"stopped no-op", Stopped, 5 * time.Second, time.Second, 5 * time.Second, false, Stopped},
		{"paused no-op", Paused, 5 * time.Second, time.Second, 5 * time.Second, false, Paused},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(tt.initialRemaining)
			timer.State = tt.initialState

			complete := timer.Tick(tt.tickAmount)

			if complete != tt.expectedComplete {
				t.Errorf("complete = %v, want %v", complete, tt.expectedComplete)
			}
			if timer.Remaining != tt.expectedRemaining {
				t.Errorf("Remaining = %v, want %v", timer.Remaining, tt.expectedRemaining)
			}
			if timer.State != tt.expectedState {
				t.Errorf("State = %v, want %v", timer.State, tt.expectedState)
			}
		})
	}
}

func TestIsRunning(t *testing.T) {
	tests := []struct {
		name     string
		state    State
		expected bool
	}{
		{"Running", Running, true},
		{"Stopped", Stopped, false},
		{"Paused", Paused, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(time.Minute)
			timer.State = tt.state
			if got := timer.IsRunning(); got != tt.expected {
				t.Errorf("IsRunning() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsComplete(t *testing.T) {
	tests := []struct {
		name      string
		remaining time.Duration
		expected  bool
	}{
		{"zero", 0, true},
		{"negative", -time.Second, true},
		{"positive", time.Second, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timer := New(time.Minute)
			timer.Remaining = tt.remaining
			if got := timer.IsComplete(); got != tt.expected {
				t.Errorf("IsComplete() = %v, want %v", got, tt.expected)
			}
		})
	}
}
