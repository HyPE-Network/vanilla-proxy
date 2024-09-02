package entity

type EntityData struct {
	TypeID    string
	RuntimeID uint64
}

type Entities map[int64]EntityData

// Init initializes the entities package.
func Init() Entities {
	return make(Entities)
}

// SetEntity sets the entity type of an entity with the specified runtime ID.
func (entities Entities) SetEntity(actorId int64, data EntityData) {
	entities[actorId] = data
}

// GetEntity returns the entity type of an entity with the specified runtime ID.
func (entities Entities) GetEntity(actorId int64) (EntityData, bool) {
	entityData, ok := entities[actorId]
	return entityData, ok
}

// GetEntityTypeID returns the entity type ID of the entity with the specified entity ID.
func (entities Entities) GetEntityTypeID(actorId int64) (string, bool) {
	entityData, ok := entities[actorId]
	if !ok {
		return "", ok
	}
	return entityData.TypeID, ok
}

// GetEntityRuntimeID returns the runtime ID of the entity with the specified entity ID.
func (entities Entities) GetEntityRuntimeID(actorId int64) (uint64, bool) {
	entityData, ok := entities[actorId]
	if !ok {
		return 0, ok
	}
	return entityData.RuntimeID, ok
}

// GetEntityFromRuntimeID returns the entity ID of the entity with the specified runtime ID.
func (entities Entities) GetEntityFromRuntimeID(runtimeID uint64) (int64, bool) {
	for actorID, entityData := range entities {
		if entityData.RuntimeID == runtimeID {
			return actorID, true
		}
	}
	return 0, false
}

// RemoveEntity removes the entity with the specified runtime ID.
func (entities Entities) RemoveEntity(actorId int64) {
	delete(entities, actorId)
}
