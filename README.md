# twitter-clone-go
## 起動方法
```shell
git clone git@github.com:kindaidai/twitter-clone-go.git

cd twitter-clone-go

docker compose up --build mysql server
```

ブラウザで、`http://loaclhost:8080` へアクセス

## 使用した主なライブラリ
- github.com/gin-gonic/gin
  - web frameworkとして使用
- gorm.io/gorm
  - ORMとして使用
- golang.org/x/crypto/bcrypt
  - passwordのハッシュ化
## 画面一覧
- TOPページ(`GET /`)
  - フォローユーザーと自分のツイートがtweetsテーブルのidで降順で表示される
  - Tweetする(`POST /tweet`)
- ログイン(`GET /login`)
  - email, passwordでログインする
  - usersテーブルのidをcookieに保持する
- サインアップ(`GET /signup`)
  - name, email, passwordでサインアップする
- 未フォローユーザー一覧(`GET /users`)
  - フォローしていないユーザーが表示される
  - フォローする(`POST /follow`)
## DB構成
![db](https://user-images.githubusercontent.com/19383278/130703273-6849287f-b65f-4089-b7c3-b652c47bf80d.png)
```sql
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(20) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` longblob NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4

CREATE TABLE `tweets` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `content` text NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_tweets_deleted_at` (`deleted_at`),
  KEY `fk_users_tweets` (`user_id`),
  CONSTRAINT `fk_users_tweets` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4

CREATE TABLE `follows` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `follower_id` bigint(20) unsigned NOT NULL,
  `followed_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_follows_deleted_at` (`deleted_at`),
  KEY `fk_follows_followed` (`followed_id`),
  KEY `fk_follows_follower` (`follower_id`),
  CONSTRAINT `fk_follows_followed` FOREIGN KEY (`followed_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_follows_follower` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
```

- usersのレコードを削除した時に、usersのidを外部キーとして持っているテーブルの外部キーには、CASCADEを設定して、usersのレコードの削除と同時に削除するようにしています。
- followsは、フォローの関係を表す中間テーブルとして作成しています。
- マイグレーションは、gormの[Auto Migration](https://gorm.io/docs/migration.html
)を使用しています。

## 最低限実装いただきたい機能の他に追加した機能
- ログイン機能
  - email, passwordによるログイン機能を追加しました。
  - emailに一意制約をつけているので、emailでユーザーをDBから取得し、passwordによる付き合わせをしています。

## やり残したこと
- cookie改ざん対策
  - ログインしたユーザーのidをそのままcookieに保持し、そのidを用いてログインユーザーを特定する処理にしているため、最低限暗号化したかったです。
- バリデーション
  - フォームから受け取った値に対して、意図しないデータを作成させないようにバリデーションをつけ、アプリケーション側でデータチェックをしたかったです。
- go buildで、バイナリからサーバーを起動させること
  - Dockerfileで`RUN go buil -o twitter-clone-go`でimageをbuild時にバイナリを作成し、`docker compose`でコンテナを起動時にバイナリを実行しようとしましたが、コンテナ起動時に作成したバイナリが見つからず、原因がまだわかっていなかったので`go run main.go`でサーバーを起動するようにしています。
- テスト
  - gormによるDB操作は最低限テストをかきたかったです。
- ディレクトリ構成の考慮
  - `main.go`に全て処理を集約させてしまっているので、ディレクトリ構成を細かくし、可読性を考慮したかったです。
  - 以下のような、GETの第二引数のfunctionの切り出しができておらず、可読性が落ちているので、ファイル切り出しもみつつ可読性をあげたかったです。
```
router.GET("/signin", func(c *gin.Context) {
        c.HTML(http.StatusOK, "signin.html", gin.H{})
})
```
