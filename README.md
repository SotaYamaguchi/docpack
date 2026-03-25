# meeting-toolkit

[![Test](https://github.com/SotaYamaguchi/meeting-toolkit/actions/workflows/test.yml/badge.svg)](https://github.com/SotaYamaguchi/meeting-toolkit/actions/workflows/test.yml)

MTG前後の資料ファイルを自動で整理するCLIツール。

## インストール

```bash
./install.sh
```

自動セットアップを選択すると、シェルを再起動するだけで使えます：

```bash
exec zsh
```

## 使い方

```bash
# プロジェクト一覧
mtg list

# MTG前の資料準備
mtg prep -project your-project

# MTG後の議事メモ整理
mtg memo -project your-project

# メールテンプレート作成（初回のみ）
mtg mail init -project your-project -type prep  # MTG前送付メール用
mtg mail init -project your-project -type memo  # MTG後送付メール用

# メールテンプレート表示
mtg mail -project your-project -type prep  # MTG前送付メール
mtg mail -project your-project -type memo  # MTG後送付メール
```

## 初期設定

### プロジェクト設定

`~/.config/mtg/config.json` を編集してプロジェクトを追加：

```json
{
  "projects": {
    "project-a": "PREFIX_A",
    "project-b": "PREFIX_B"
  },
  "mail_templates": {
    "project-a": {
      "prep": "templates/project-a-prep.txt",
      "memo": "templates/project-a-memo.txt"
    }
  }
}
```

### メールテンプレート設定

#### 方法1: コマンドで作成（推奨）

```bash
# テンプレートファイルを自動作成
mtg mail init -project your-project -type prep

# ✓ テンプレートファイルを作成しました: ~/.config/mtg/templates/your-project-prep.txt
# ✓ config.jsonを更新しました
#
# テンプレートを編集してください:
#   vim ~/.config/mtg/templates/your-project-prep.txt
```

作成されたテンプレートを編集：

```
To: customer@example.com, another@example.com
Cc: team@example.com
Subject: 【プロジェクトA】MTG資料送付 {{DATE}}

お世話になっております。

本日のMTG資料を送付いたします。

送付資料：
- 資料_{{DATE}}.pdf

ご確認のほど、よろしくお願いいたします。
```

#### 方法2: 手動で作成

`~/.config/mtg/templates/` にテンプレートファイルを直接作成し、`config.json` の `mail_templates` セクションに追加することもできます。

**特徴：**
- メーラーからのコピペがそのまま使える
- 改行や箇条書きもそのまま保持される
- To/Cc/Bccはカンマ区切りで複数指定可能
- `{{DATE}}` は実行日の日付（YYYYMMDD形式）に自動置換

## 詳細

- タブ補完対応（zsh）
- ヘルプ: `mtg help`
- アンインストール: `cd mtg && make uninstall`
- 開発者向け: [CONTRIBUTING.md](CONTRIBUTING.md)
