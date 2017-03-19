# reicon
command line tool for changing twitter profile icon

## なにこれ

twitterのプロフィール画像をコマンドラインからランダムで(<-ココ重要)変更します

## 使い方

ターミナルからコマンドで実行できます。

- **config** コンフィグファイルを指定します
- **image** 画像をパスで指定します。Globでランダムに出す画像を指定できます。;区切りで更に画像を上に乗せれます。
- **color** ベースの色です。#ffffff or #ffffffff の形式です。カンマ区切りで複数指定できます。(省略時透明)

初回時に以下の形式で実行すると設定ファイルが作成されます。
二回目でtwitterの承認のURLが表示されるので、そのURLをブラウザで表示してピンをもらってください。
それを"PIN:"のあとで入力すると承認完了です。
```
reicon -config /home/user/reicon/config.json
```

ex. config.jsonの設定でベースカラー赤or青でimg配下の画像をランダムで合成したものをアイコンにします
```
reicon -config /home/user/reicon/config.json -image "/home/user/img/*" -color "#ff0000, #0000ff"
```

ex. config.jsonの設定でbg配下の画像を下地にimg配下の画像をランダムで合成したものをアイコンにします
```
reicon -config /home/user/reicon/config.json -image "/home/user/bg/*;/home/user/bg/*"
```

主な用途は crontab に指定してtwitterのアイコンを定期的に変えることです。

### 色がわからない

webカラーコード等でググってください

### 画像サイズが均一でない場合は？

正方形にピッタリ入るサイズではみ出した部分は切り取られます

### 画像の最大サイズは？

400x400ピクセルにリサイズされます

## ダウンロード

| 環境                      | URL |
|:-------------------------:|:-:|
| Windows                   | https://github.com/Mushus/reicon/raw/master/build/windows-amd64/reicon.exe |
| Linux                     | https://github.com/Mushus/reicon/raw/master/build/linux-amd64/reicon |
| MacOS                     | https://github.com/Mushus/reicon/raw/master/build/darwin-amd64/reicon |
| Linux/ARM6(rasberry Pi 1) | https://github.com/Mushus/reicon/raw/master/build/linux-arm6/reicon |
| ARM7                      | https://github.com/Mushus/reicon/raw/master/build/linux-arm7/reicon |

## Licence

MIT
