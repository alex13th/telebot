package fsm

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var states = map[string][]State{
	"100": {{ChatId: "100", State: "state1", Key: "key1"}},
	"101": {{ChatId: "101", State: "state11", Key: "key2"}, {ChatId: "101", MessageId: 111, State: "state11", Key: "key1"}},
	"102": {{ChatId: "102", State: "state12", Key: "key1"}},
}

func TestState_Parse(t *testing.T) {
	tests := []struct {
		name    string
		state   State
		data    string
		want    State
		wantErr bool
	}{
		{
			name:  "Minimal",
			state: State{Separator: "_"},
			data:  "pr1_state1_action1",
			want:  State{Prefix: "pr1", Separator: "_", State: "state1", Action: "action1"},
		},
		{
			name:  "With Data",
			state: State{Separator: "_"},
			data:  "pr2_state1_action1_key1",
			want:  State{Prefix: "pr2", Separator: "_", State: "state1", Key: "key1", Action: "action1"},
		},
		{
			name:  "With Value",
			state: State{Separator: "-"},
			data:  "pr3-state1-action1-key2-value-contains-separator",
			want:  State{Prefix: "pr3", Separator: "-", State: "state1", Key: "key2", Action: "action1", Value: "value-contains-separator"},
		},
		{
			name:    "Incorrect parts error",
			state:   State{Separator: "-"},
			data:    "pr3-state1",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.state.Parse(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("State.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("State.Parse() difference: %s", diff)
			}
		})
	}
}

func TestState_String(t *testing.T) {
	tests := []struct {
		name  string
		state State
		want  string
	}{
		{
			name:  "Minimal",
			state: State{Prefix: "pr", Separator: "_", State: "state1"},
			want:  "pr_state1_state1",
		},
		{
			name:  "With action",
			state: State{Prefix: "pr", Separator: "_", State: "state1", Action: "action1"},
			want:  "pr_state1_action1",
		},
		{
			name:  "With Data",
			state: State{Prefix: "pr", Separator: "_", State: "state1", Key: "key1", Action: "action1"},
			want:  "pr_state1_action1_key1",
		},
		{
			name:  "With Value",
			state: State{Prefix: "pr", Separator: "_", State: "state1", Key: "key2", Action: "action1", Value: "value_contains_separator"},
			want:  "pr_state1_action1_key2_value_contains_separator",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.String(); got != tt.want {
				t.Errorf("State.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewState(t *testing.T) {
	t.Run("New state", func(t *testing.T) {
		if got := NewState(); got.Separator != "_" {
			t.Errorf("NewState() = %v, want %v", got.Separator, "_")
		}
	})
}

func TestNewMemoryStateRepository(t *testing.T) {
	want := MemoryStateRepository{chatStates: make(map[string][]State)}
	t.Run("NewMemoryStateRepository", func(t *testing.T) {
		got := NewMemoryStateRepository()
		if diff := cmp.Diff(got.chatStates, want.chatStates); diff != "" {
			t.Errorf("NewMemoryStateRepository() difference: %v", diff)
		}
	})
}

func TestMemoryStateRepository_Get(t *testing.T) {
	tests := []struct {
		name    string
		chatId  string
		wantSt  []State
		wantErr bool
	}{
		{
			name:   "Several states",
			chatId: "101",
			wantSt: []State{{ChatId: "101", State: "state11", Key: "key2"}, {ChatId: "101", MessageId: 111, State: "state11", Key: "key1"}},
		},
		{
			name:   "One state present",
			chatId: "102",
			wantSt: []State{{ChatId: "102", State: "state12", Key: "key1"}},
		},
		{
			name:    "State not found",
			chatId:  "1000",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := &MemoryStateRepository{chatStates: states}
			gotSt, err := rep.Get(tt.chatId)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStateRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotSt, tt.wantSt); diff != "" {
				t.Errorf("MemoryStateRepository.Get() difference: %v", diff)
			}
		})
	}
}

func TestMemoryStateRepository_GetByMessage(t *testing.T) {
	type args struct {
		chatId    string
		messageId int
	}
	tests := []struct {
		name    string
		args    args
		want    State
		wantErr bool
	}{
		{
			name: "Several states",
			args: args{chatId: "101", messageId: 111},
			want: State{ChatId: "101", MessageId: 111, State: "state11", Key: "key1"},
		},
		{
			name:    "Message state not in chat",
			args:    args{chatId: "101", messageId: 222},
			wantErr: true,
		},
		{
			name:    "Chat state not found",
			args:    args{chatId: "1000", messageId: 111},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := &MemoryStateRepository{chatStates: states}
			got, err := rep.GetByMessage(tt.args.chatId, tt.args.messageId)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStateRepository.GetByMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("MemoryStateRepository.GetByMessage() difference: %v", diff)
			}
		})
	}
}

func TestMemoryStateRepository_GetByKey(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantSlist []State
		wantErr   bool
	}{
		{
			name: "Several states",
			key:  "key1",
			wantSlist: []State{
				{ChatId: "100", State: "state1", Key: "key1"},
				{ChatId: "101", MessageId: 111, State: "state11", Key: "key1"},
				{ChatId: "102", State: "state12", Key: "key1"},
			},
		},
		{
			name:      "One state",
			key:       "key2",
			wantSlist: []State{{ChatId: "101", State: "state11", Key: "key2"}},
		},
		{
			name:    "States not found",
			key:     "key10",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := &MemoryStateRepository{chatStates: states}
			gotSlist, err := rep.GetByKey(tt.key)
			sort.Slice(gotSlist, func(i, j int) bool { return gotSlist[i].ChatId < gotSlist[j].ChatId })
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStateRepository.GetByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(gotSlist, tt.wantSlist); diff != "" {
				t.Errorf("MemoryStateRepository.GetByKey() = difference %v", diff)
			}
		})
	}
}

func TestMemoryStateRepository_Set(t *testing.T) {
	tStates := make(map[string][]State, len(states))
	for k, v := range states {
		tStates[k] = v
	}
	tests := []struct {
		name       string
		chatStates map[string][]State
		st         State
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "Update state",
			chatStates: tStates,
			st:         State{ChatId: "101", MessageId: 111, State: "state11", Key: "key1"},
			wantLen:    len(tStates),
		},
		{
			name:       "New state",
			chatStates: tStates,
			st:         State{ChatId: "1101", MessageId: 111, State: "state11", Key: "key1"},
			wantLen:    len(tStates) + 1,
		},
		{
			name:    "Empty repository",
			st:      State{ChatId: "1101", MessageId: 111, State: "state11", Key: "key1"},
			wantLen: 1,
		},
		{
			name:       "Empty ChatId error",
			chatStates: tStates,
			st:         State{MessageId: 111, State: "state11", Key: "key1"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := &MemoryStateRepository{chatStates: tt.chatStates}

			err := rep.Set(tt.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStateRepository.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			got, err := rep.GetByMessage(tt.st.ChatId, tt.st.MessageId)
			if err != nil {
				t.Errorf("MemoryStateRepository.Set() get state after setting error = %v", err)
			}
			if diff := cmp.Diff(got, tt.st); diff != "" {
				t.Errorf("MemoryStateRepository.Set() = difference %v", diff)
			}
			if len(rep.chatStates) != tt.wantLen {
				t.Errorf("MemoryStateRepository.Set() chats count wanted %v, but  %v", len(tStates), tt.wantLen)
			}
		})
	}
}

func TestMemoryStateRepository_Clear(t *testing.T) {
	tStates := make(map[string][]State, len(states))
	for k, v := range states {
		tStates[k] = v
	}
	tests := []struct {
		name        string
		chatStates  map[string][]State
		st          State
		wantLen     int
		wantChatLen int
		wantErr     bool
	}{
		{
			name:        "Clear message state",
			chatStates:  tStates,
			st:          State{ChatId: "101", MessageId: 111, State: "state11"},
			wantLen:     len(tStates),
			wantChatLen: len(tStates["101"]) - 1,
		},
		{
			name:        "Clear chat state",
			chatStates:  tStates,
			st:          State{ChatId: "100", State: "state1"},
			wantChatLen: len(tStates["100"]) - 1,
			wantLen:     len(tStates),
		},
		{
			name:       "Clear all chat states",
			chatStates: tStates,
			st:         State{ChatId: "102"},
			wantLen:    len(tStates) - 1,
		},
		{
			name:    "Empty repository",
			st:      State{ChatId: "102"},
			wantLen: 0,
		},
		{
			name:       "Empty ChatId error",
			chatStates: tStates,
			st:         State{MessageId: 111, State: "state11"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := &MemoryStateRepository{chatStates: tt.chatStates}

			err := rep.Clear(tt.st)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemoryStateRepository.Clear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			_, err = rep.GetByMessage(tt.st.ChatId, tt.st.MessageId)
			if err == nil {
				t.Error("MemoryStateRepository.Clear() get state after clear must raise error")
			}
			if len(rep.chatStates) != tt.wantLen {
				t.Errorf("MemoryStateRepository.Clear() chats count wanted %v, but  %v",
					len(rep.chatStates), tt.wantLen)
			}
			if len(rep.chatStates[tt.st.ChatId]) != tt.wantChatLen {
				t.Errorf("MemoryStateRepository.Clear() chats count wanted %v, but  %v",
					len(rep.chatStates[tt.st.ChatId]), tt.wantChatLen)
			}
		})
	}
}
