package entity

import "sync"

type EntityData struct {
	TypeID    string
	RuntimeID uint64
}

type Entities struct {
	IdToData map[int64]EntityData
	mu       sync.RWMutex
}

// Init initializes the entities package.
func Init() *Entities {
	return &Entities{
		IdToData: make(map[int64]EntityData),
	}
}

// SetEntity sets the entity type of an entity with the specified runtime ID.
func (entities *Entities) SetEntity(actorId int64, data EntityData) {
	entities.mu.Lock()
	defer entities.mu.Unlock()
	entities.IdToData[actorId] = data
}

// GetEntity returns the entity type of an entity with the specified runtime ID.
func (entities *Entities) GetEntity(actorId int64) (EntityData, bool) {
	entities.mu.RLock()
	defer entities.mu.RUnlock()
	entityData, ok := entities.IdToData[actorId]
	return entityData, ok
}

// GetEntityTypeID returns the entity type ID of the entity with the specified entity ID.
func (entities *Entities) GetEntityTypeID(actorId int64) (string, bool) {
	entities.mu.RLock()
	defer entities.mu.RUnlock()

	entityData, ok := entities.IdToData[actorId]
	if !ok {
		return "", ok
	}
	return entityData.TypeID, ok
}

// GetEntityRuntimeID returns the runtime ID of the entity with the specified entity ID.
func (entities *Entities) GetEntityRuntimeID(actorId int64) (uint64, bool) {
	entities.mu.RLock()
	defer entities.mu.RUnlock()

	entityData, ok := entities.IdToData[actorId]
	if !ok {
		return 0, ok
	}
	return entityData.RuntimeID, ok
}

// GetEntityFromRuntimeID returns the entity ID of the entity with the specified runtime ID.
func (entities *Entities) GetEntityFromRuntimeID(runtimeID uint64) (int64, bool) {
	entities.mu.RLock()
	defer entities.mu.RUnlock()

	for actorID, entityData := range entities.IdToData {
		if entityData.RuntimeID == runtimeID {
			return actorID, true
		}
	}
	return 0, false
}

// RemoveEntity removes the entity with the specified runtime ID.
func (entities *Entities) RemoveEntity(actorId int64) {
	entities.mu.Lock()
	defer entities.mu.Unlock()

	delete(entities.IdToData, actorId)
}
