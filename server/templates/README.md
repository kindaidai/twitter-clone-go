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

## TODO
- cookie改ざん対策
  - ログインしたユーザーのidをそのままcookieに保持し、そのidを用いてログインユーザーを特定する処理にしているため、最低限暗号化。
- バリデーション
  - フォームから受け取った値に対して、意図しないデータを作成させないようにバリデーションをつけ、アプリケーション側でデータチェック。
- go buildで、バイナリからサーバーを起動させること
  - Dockerfileで`RUN go buil -o twitter-clone-go`でimageをbuild時にバイナリを作成し、`docker compose`でコンテナを起動時にバイナリを実行しようとしましたが、コンテナ起動時に作成したバイナリが見つからず、原因がまだわかっていなかったので`go run main.go`でサーバーを起動。
- テスト
  - gormによるDB操作は最低限テスト。
- ディレクトリ構成の考慮
  - `main.go`に全て処理を集約させてしまっているので、ディレクトリ構成を細かくし、可読性を考慮する。
  - 以下のような、GETの第二引数のfunctionの切り出しができていない。
```
router.GET("/signin", func(c *gin.Context) {
        c.HTML(http.StatusOK, "signin.html", gin.H{})
})
```