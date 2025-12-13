package mysql

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/kinpatsu-everyone/backend-template/config"
)

var (
	db   *sql.DB
	once sync.Once
)

// PoolConfig はコネクションプールの設定を表します
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultPoolConfig は本番環境向けのデフォルト設定を返します
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

const tlsConfigName = "tidb"

// InitDB はデータベース接続を初期化します
// アプリケーション起動時に一度だけ呼び出してください
func InitDB(ctx context.Context, poolCfg PoolConfig) error {
	var initErr error
	once.Do(func() {
		cfg := mysql.Config{
			User:                 config.MySQLUser,
			Passwd:               config.MySQLPassword,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%s", config.MySQLHost, config.MySQLPort),
			DBName:               config.MySQLDatabase,
			ParseTime:            true,
			Loc:                  time.UTC,
			Collation:            "utf8mb4_unicode_ci",
			AllowNativePasswords: true,
		}

		// クラウド環境の場合はTLS設定を有効化（TiDB Serverless向け）
		if config.IsCloud() {
			if err := mysql.RegisterTLSConfig(tlsConfigName, &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: config.MySQLHost,
			}); err != nil {
				initErr = fmt.Errorf("failed to register TLS config: %w", err)
				return
			}
			cfg.TLSConfig = tlsConfigName
		}

		conn, err := sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			initErr = fmt.Errorf("failed to open database connection: %w", err)
			return
		}

		// コネクションプール設定
		conn.SetMaxOpenConns(poolCfg.MaxOpenConns)
		conn.SetMaxIdleConns(poolCfg.MaxIdleConns)
		conn.SetConnMaxLifetime(poolCfg.ConnMaxLifetime)
		conn.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)

		// 接続確認
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := conn.PingContext(pingCtx); err != nil {
			conn.Close()
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		db = conn
	})

	return initErr
}

// GetDB はデータベース接続を返します
// InitDB が先に呼ばれている必要があります
func GetDB() *sql.DB {
	if db == nil {
		panic("mysql: database connection not initialized. Call InitDB first.")
	}
	return db
}

// GetQueries はsqlcで生成されたQueriesを返します
func GetQueries() *Queries {
	return New(GetDB())
}

// Close はデータベース接続を閉じます
// アプリケーション終了時に呼び出してください
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// HealthCheck はデータベースの接続状態を確認します
func HealthCheck(ctx context.Context) error {
	if db == nil {
		return fmt.Errorf("mysql: database connection not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("mysql: health check failed: %w", err)
	}
	return nil
}

// Stats はコネクションプールの統計情報を返します
func Stats() sql.DBStats {
	if db == nil {
		return sql.DBStats{}
	}
	return db.Stats()
}

// WithTx はトランザクション内で関数を実行します
// エラーが発生した場合は自動的にロールバックされます
func WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := GetDB().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// WithQueriesTx はトランザクション内でQueriesを使って処理を実行します
func WithQueriesTx(ctx context.Context, fn func(q *Queries) error) error {
	return WithTx(ctx, func(tx *sql.Tx) error {
		q := New(tx)
		return fn(q)
	})
}
