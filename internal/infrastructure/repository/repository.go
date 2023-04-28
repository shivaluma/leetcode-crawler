package repository

type Repository struct {
	Lastupdatedtimestamp LastUpdatedTimestampRepositoryInterface
}

func NewRepository(lut LastUpdatedTimestampRepositoryInterface) *Repository {
	return &Repository{
		Lastupdatedtimestamp: lut,
	}
}
