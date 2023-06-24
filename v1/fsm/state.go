package fsm

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrStateNotFound    = errors.New("the state was not found in the repository")
	ErrFailedToAddState = errors.New("failed to add the state to the repository")
	ErrUpdateState      = errors.New("failed to update the state in the repository")
)

func NewState() State {
	return State{Separator: "_"}
}

type State struct {
	Action    string `json:"action"`
	ChatId    string `json:"chat_id"`
	Key       string `json:"key"`
	MessageId int    `json:"message_id"`
	Prefix    string `json:"prefix"`
	Separator string `json:"separator"`
	State     string `json:"state"`
	Value     string `json:"value"`
}

func (st State) String() string {
	if st.Action == "" {
		st.Action = st.State
	}
	slist := []string{st.Prefix, st.State, st.Action}
	if st.Key != "" {
		slist = append(slist, st.Key)
	}
	if st.Value != "" {
		slist = append(slist, st.Value)
	}
	return strings.Join(slist, st.Separator)
}

func (st State) Parse(data string) (State, error) {
	state := st
	splitedData := strings.Split(data, st.Separator)
	if len(splitedData) < 3 {
		return State{},
			fmt.Errorf("incorrect button data format, data (%s) must has at least 3 parts, but has %d ",
				splitedData, len(splitedData))
	}
	state.Prefix = splitedData[0]
	state.State = splitedData[1]
	state.Action = splitedData[2]
	if len(splitedData) > 3 {
		state.Key = splitedData[3]
	}
	if len(splitedData) > 4 {
		state.Value = strings.Join(splitedData[4:], st.Separator)
	}
	return state, nil
}

func NewMemoryStateRepository() MemoryStateRepository {
	return MemoryStateRepository{chatStates: make(map[string][]State)}
}

type MemoryStateRepository struct {
	chatStates map[string][]State
	sync.Mutex
}

func (rep *MemoryStateRepository) Get(chatId string) (st []State, err error) {
	if states, ok := rep.chatStates[chatId]; ok {
		return states, nil
	}

	return nil, ErrStateNotFound
}

func (rep *MemoryStateRepository) GetByMessage(chatId string, messageId int) (State, error) {
	states, err := rep.Get(chatId)
	if err != nil {
		return State{}, err
	}
	for _, s := range states {
		if s.MessageId == messageId {
			return s, nil
		}
	}
	return State{}, ErrStateNotFound
}

func (rep *MemoryStateRepository) GetByKey(key string) (slist []State, err error) {
	rep.Mutex.Lock()
	for _, ss := range rep.chatStates {
		for _, s := range ss {
			if s.Key == key {
				slist = append(slist, s)
			}
		}
	}
	rep.Mutex.Unlock()
	if len(slist) == 0 {
		return nil, ErrStateNotFound
	}
	return
}

func (rep *MemoryStateRepository) Set(s State) error {
	if s.ChatId == "" {
		return fmt.Errorf("State ChatId can't be empty, state: %v", s)
	}

	rep.Lock()
	defer rep.Unlock()
	if rep.chatStates == nil {
		rep.chatStates = make(map[string][]State)
	}
	rep.chatStates[s.ChatId] = []State{s}
	return nil
}

func (rep *MemoryStateRepository) Clear(st State) error {
	if st.ChatId == "" {
		return fmt.Errorf("State ChatId can't be empty, state: %v", st)
	}

	rep.Lock()
	defer rep.Unlock()

	if rep.chatStates == nil {
		rep.chatStates = make(map[string][]State)
	} else {
		if st.MessageId == 0 && st.State == "" {
			delete(rep.chatStates, st.ChatId)
			return nil
		}
		for k, v := range rep.chatStates {
			for i, s := range v {
				if st.MessageId != 0 && s.ChatId == st.ChatId && s.MessageId == st.MessageId {
					rep.chatStates[k] = append(v[:i], v[i+1:]...)
				}
				if st.MessageId == 0 && s.ChatId == st.ChatId && s.State == st.State {
					rep.chatStates[k] = append(v[:i], v[i+1:]...)
				}
			}
		}
	}
	return nil
}
