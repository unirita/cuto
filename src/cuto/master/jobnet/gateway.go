package jobnet

import "fmt"

// ジョブネット内のフロー分岐・集約を表す構造体
type Gateway struct {
	id    string
	Nexts []Element
}

// Gateway構造体のコンストラクタ関数
func NewGateway(id string) *Gateway {
	gateway := new(Gateway)
	gateway.id = id
	return gateway
}

// IDを取得する
func (g *Gateway) ID() string {
	return g.id
}

// ノードタイプを取得する
//
// return : ノードタイプ
func (g *Gateway) Type() elementType {
	return ELM_GW
}

// 後続エレメントの追加を行う。
func (g *Gateway) AddNext(e Element) error {
	g.Nexts = append(g.Nexts, e)
	return nil
}

// 後続エレメントの有無を調べる。
func (g *Gateway) HasNext() bool {
	return len(g.Nexts) > 0
}

// 後続のノード数に応じて、ゲートウェイの処理を行う。
//
// 後続ノードが存在しない場合：何も行わない。次の実行ノードとしてnilを返す。
//
// 後続ノードが1つだけの場合：何も行わない。次の実行ノードとして唯一の後続ノードを返す。
//
// 後続ノードが2つ以上の場合：各後続ノードを先頭としたPath構造体を生成し、並列実行する。次の実行ノードとして結合ゲートウェイを返す。
//
// return : 次の実行ノード
//
// return : エラー情報
func (g *Gateway) Execute() (Element, error) {
	nextCnt := len(g.Nexts)

	switch nextCnt {
	case 0:
		return nil, nil
	case 1:
		return g.Nexts[0], nil
	default:
		paths := make([]*Path, nextCnt)
		done := make(chan struct{}, nextCnt)
		for i := 0; i < nextCnt; i++ {
			paths[i] = NewPath(g.Nexts[i])
			go paths[i].Run(done)
		}
		for i := 0; i < nextCnt; i++ {
			<-done
		}
		for i := 0; i < nextCnt; i++ {
			if paths[i].Err != nil {
				err := fmt.Errorf("Error occured in branch: %s", paths[i].Err)
				return nil, err
			}
			if paths[i].Goal != paths[0].Goal {
				err := fmt.Errorf("Branch is combined by more than one gateway.")
				return nil, err
			}
		}
		return paths[0].Goal, nil
	}
}
