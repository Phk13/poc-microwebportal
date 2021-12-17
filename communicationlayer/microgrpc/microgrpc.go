package microgrpc

import (
	"context"
	"errors"

	"github.com/phk13/poc-micro/databaselayer"
)

type MicroGrpcServer struct {
	dbHandler databaselayer.DinoDBHandler
}

func NewMicroGrpcServer(dbtype uint8, connstring string) (*MicroGrpcServer, error) {
	handler, err := databaselayer.GetDatabaseHandler(dbtype, connstring)
	if err != nil {
		return nil, errors.New("could not create a database handler object")
	}
	return &MicroGrpcServer{
		dbHandler: handler,
	}, nil
}

func (server *MicroGrpcServer) GetAnimal(ctx context.Context, r *Request) (*Animal, error) {
	animal, err := server.dbHandler.GetDinoByNickname(r.GetNickname())
	return convertToDinoGRPCAnimal(animal), err
}

func (server *MicroGrpcServer) GetAllAnimals(req *Request, stream MicroService_GetAllAnimalsServer) error {
	animals, err := server.dbHandler.GetAvailableDinos()
	if err != nil {
		return err
	}
	for _, animal := range animals {
		grpcAnimal := convertToDinoGRPCAnimal(animal)
		if err := stream.Send(grpcAnimal); err != nil {
			return err
		}
	}
	return nil
}

func convertToDinoGRPCAnimal(animal databaselayer.Animal) *Animal {
	return &Animal{
		Id:         int32(animal.ID),
		AnimalType: animal.AnimalType,
		Nickname:   animal.Nickname,
		Zone:       int32(animal.Zone),
		Age:        int32(animal.Age),
	}
}
