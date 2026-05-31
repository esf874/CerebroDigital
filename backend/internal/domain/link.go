package domain

// Link representa una conexión entre dos notas del grafo.
type Link struct {
	ID           string `bson:"_id,omitempty"`
	Alias        string `bson:"alias,omitempty"`
	OriginNoteId string `bson:"origin_note_id"`
	DestNoteId   string `bson:"dest_note_id"`
}

func NewLink(id, originID, destID, alias, destTitle string) (*Link, error) {
	if id == "" || originID == "" || destID == "" {
		return nil, ErrInvalidID
	}

	// Comprobar que no se crea un enlace a la misma nota
	if originID == destID {
		return nil, ErrSameIds
	}

	// Si no hay alias, uso de titulo nota destino por defecto
	finalAlias := alias
	if finalAlias == "" {
		finalAlias = destTitle
	}

	return &Link{
		ID:           id,
		OriginNoteId: originID,
		DestNoteId:   destID,
		Alias:        finalAlias,
	}, nil
}

// No se permite el cambio de nota origen o destino, si se quiere cambiar,
// borrar y crear un nuevo enlace.
