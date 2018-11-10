
# アプリ仕様
- 質問(Question),回答(Answer)のセットを、指定位置にボタンとして表示。
- ボタンにより質問、回答が交互にボタンラベルとなる。
- ラベルが質問の場合、未回答、回答の場合、回答としてカウントされる。
- 保存のタイミング（回答がされる度）に未回答、回答がOKcout,NGcountとして計算される。
- 保存のタイミングで利用者(User)の各質問ページ情報(QuestionPage)に現在の状態(Status)に回答の質問(QuestionInfo)のID（Id)が登録される。
- 質問ページ(QuestionPage)が読み込まれる時は、QuestionPageの状態(Status)にそって書く質問を表示する。

## TODO
- [x] ボタン表示の切り替え、簡易位置設定
- [x] 既存Status表示 

- [x] 保存タイミングのStatsuオブジェクト変更
- [x] Statusオブジェクトの永続化
- [x] 保存後の再読み込み
- [x] 初期読み込み表示　永続化からの読み込み
- [ ] top pageの作成
- [ ] User情報の作成,映像化

## 
- POST request.ParsePost で hidden inputのstatus-infoに id-statusOK, id-statusNGで各質問の現状を保持

