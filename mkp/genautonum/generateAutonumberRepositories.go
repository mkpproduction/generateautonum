package genautonum

type generateAutonumberRepository struct {
	RepoDB Repository
}

func NewGenerateAutonumberRepository(repoDB Repository) generateAutonumberRepository {
	return generateAutonumberRepository{
		RepoDB: repoDB,
	}
}

type GenerateAutonumberRepository interface {
	GenerateAutonumber(p string, v string) (string, error)
}

// GenerateAutonumber
func (ctx generateAutonumberRepository) GenerateAutonumber(p string, v string) (string, error) {
	var autonumber string

	err := ctx.RepoDB.DB.QueryRow("SELECT fs_gen_autonum($1, $2)", p, v).Scan(&autonumber)
	if err != nil {
		return "", err
	}

	return autonumber, nil
}