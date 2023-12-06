#!/bin/sh

if ! (type task >/dev/null 2>&1); then
  echo "下記公式ドキュメントを参考に、taskをインストールした後、もう一度init.shを実行してください"
  echo "https://taskfile.dev/ja-JP/installation/"
fi

task init
