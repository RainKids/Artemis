package repository

import (
	"admin/pkg/database/es"
	"admin/pkg/database/mongo"
	"admin/pkg/database/postgres"
	"admin/pkg/database/redis"
	"context"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type repository struct {
	db       *gorm.DB
	rdb      *redis.RedisDB
	es       *es.Client
	mongo    *mongo.MongoDB
	hello    HelloRepository
	api      ApiRepository
	casbin   CasbinRepository
	dept     DeptRepository
	dict     DictRepository
	menu     MenuRepository
	post     PostRepository
	role     RoleRepository
	user     UserRepository
	migrants []Migrant
}

func NewRepository(log *zap.Logger, db *postgres.DB, rdb *redis.RedisDB, es *es.Client, mongo *mongo.MongoDB, enforcer *casbin.SyncedEnforcer) Repository {
	r := &repository{
		db:     db.Postgres,
		rdb:    rdb,
		es:     es,
		mongo:  mongo,
		hello:  newHelloRepository(log),
		api:    newApiRepository(log, db, rdb),
		casbin: newCasbinRepository(log, db),
		dept:   newDeptRepository(log, db, rdb),
		dict:   newDictRepository(log, db, rdb),
		menu:   newMenuRepository(log, db, rdb),
		post:   newPostRepository(log, db, rdb),
		role:   newRoleRepository(log, db, rdb, enforcer),
		user:   newUserRepository(log, db, rdb),
	}
	r.migrants = getMigrants(
		r.Api(),
		r.Casbin(),
		r.Dept(),
		r.Dict(),
		r.Menu(),
		r.Post(),
		r.Role(),
		r.User(),
	)
	err := r.Init()
	if err != nil {
		log.Error("DB Table Migrate Error", zap.Error(err))
	}
	return r
}

func getMigrants(objs ...interface{}) []Migrant {
	var migrants []Migrant
	for _, obj := range objs {
		if m, ok := obj.(Migrant); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}

	if r.rdb != nil {
		if err := r.rdb.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	if err = db.PingContext(ctx); err != nil {
		return err
	}

	if r.rdb == nil {
		return nil
	}
	if _, err := r.rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func (r *repository) Migrate() error {
	for _, m := range r.migrants {
		if err := m.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Hello() HelloRepository {
	return r.hello
}

func (r *repository) Api() ApiRepository {
	return r.api
}

func (r *repository) Casbin() CasbinRepository {
	return r.casbin
}

func (r *repository) Dept() DeptRepository {
	return r.dept
}

func (r *repository) Dict() DictRepository {
	return r.dict
}

func (r *repository) Menu() MenuRepository {
	return r.menu
}

func (r *repository) Post() PostRepository {
	return r.post
}

func (r *repository) Role() RoleRepository {
	return r.role
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Init() error {
	return r.Migrate()
}
