package storage

import (
	"fmt"

	cm "github.com/paradigm-network/paradigm/common"
	"github.com/paradigm-network/paradigm/errors"
	"github.com/paradigm-network/paradigm/types"
)

type Key struct {
	x string
	y string
}

func (k Key) ToString() string {
	return fmt.Sprintf("{%s, %s}", k.x, k.y)
}

type ParentRoundInfo struct {
	Round  int
	IsRoot bool
}

func NewBaseParentRoundInfo() ParentRoundInfo {
	return ParentRoundInfo{
		Round:  -1,
		IsRoot: false,
	}
}

func getValues(mapping map[string]int) []int {
	keys := make([]int, len(mapping))
	i := 0
	for _, id := range mapping {
		keys[i] = id
		i++
	}
	return keys
}

//------------------------------------------------------------------------------

type ParticipantEventsCache struct {
	participants map[string]int
	rim          *cm.RollingIndexMap
}

func NewParticipantEventsCache(size int, participants map[string]int) *ParticipantEventsCache {
	return &ParticipantEventsCache{
		participants: participants,
		rim:          cm.NewRollingIndexMap(size, getValues(participants)),
	}
}

func (pec *ParticipantEventsCache) participantID(participant string) (int, error) {
	id, ok := pec.participants[participant]
	if !ok {
		return -1, errors.NewStoreErr(errors.UnknownParticipant, participant)
	}
	return id, nil
}

//return participant events with index > skip
func (pec *ParticipantEventsCache) Get(participant string, skipIndex int) ([]string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return []string{}, err
	}

	pe, err := pec.rim.Get(id, skipIndex)
	if err != nil {
		return []string{}, err
	}

	res := make([]string, len(pe))
	for k := 0; k < len(pe); k++ {
		res[k] = pe[k].(string)
	}
	return res, nil
}

func (pec *ParticipantEventsCache) GetItem(participant string, index int) (string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return "", err
	}

	item, err := pec.rim.GetItem(id, index)
	if err != nil {
		return "", err
	}
	return item.(string), nil
}

func (pec *ParticipantEventsCache) GetLast(participant string) (string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return "", err
	}

	last, err := pec.rim.GetLast(id)
	if err != nil {
		return "", err
	}
	return last.(string), nil
}

func (pec *ParticipantEventsCache) Set(participant string, hash string, index int) error {
	id, err := pec.participantID(participant)
	if err != nil {
		return err
	}
	return pec.rim.Set(id, hash, index)
}

//returns [participant id] => lastKnownIndex
func (pec *ParticipantEventsCache) Known() map[int]int {
	return pec.rim.Known()
}

func (pec *ParticipantEventsCache) Reset() error {
	return pec.rim.Reset()
}

//------------------------------------------------------------------------------

type ParticipantBlockSignaturesCache struct {
	participants map[string]int
	rim          *cm.RollingIndexMap
}

func NewParticipantBlockSignaturesCache(size int, participants map[string]int) *ParticipantBlockSignaturesCache {
	return &ParticipantBlockSignaturesCache{
		participants: participants,
		rim:          cm.NewRollingIndexMap(size, getValues(participants)),
	}
}

func (psc *ParticipantBlockSignaturesCache) participantID(participant string) (int, error) {
	id, ok := psc.participants[participant]
	if !ok {
		return -1, errors.NewStoreErr(errors.UnknownParticipant, participant)
	}
	return id, nil
}

//return participant BlockSignatures where index > skip
func (psc *ParticipantBlockSignaturesCache) Get(participant string, skipIndex int) ([]types.BlockSignature, error) {
	id, err := psc.participantID(participant)
	if err != nil {
		return []types.BlockSignature{}, err
	}

	ps, err := psc.rim.Get(id, skipIndex)
	if err != nil {
		return []types.BlockSignature{}, err
	}

	res := make([]types.BlockSignature, len(ps))
	for k := 0; k < len(ps); k++ {
		res[k] = ps[k].(types.BlockSignature)
	}
	return res, nil
}

func (psc *ParticipantBlockSignaturesCache) GetItem(participant string, index int) (types.BlockSignature, error) {
	id, err := psc.participantID(participant)
	if err != nil {
		return types.BlockSignature{}, err
	}

	item, err := psc.rim.GetItem(id, index)
	if err != nil {
		return types.BlockSignature{}, err
	}
	return item.(types.BlockSignature), nil
}

func (psc *ParticipantBlockSignaturesCache) GetLast(participant string) (types.BlockSignature, error) {
	last, err := psc.rim.GetLast(psc.participants[participant])
	if err != nil {
		return types.BlockSignature{}, err
	}
	return last.(types.BlockSignature), nil
}

func (psc *ParticipantBlockSignaturesCache) Set(participant string, sig types.BlockSignature) error {
	id, err := psc.participantID(participant)
	if err != nil {
		return err
	}

	return psc.rim.Set(id, sig, sig.Index)
}

//returns [participant id] => last BlockSignature Index
func (psc *ParticipantBlockSignaturesCache) Known() map[int]int {
	return psc.rim.Known()
}

func (psc *ParticipantBlockSignaturesCache) Reset() error {
	return psc.rim.Reset()
}
