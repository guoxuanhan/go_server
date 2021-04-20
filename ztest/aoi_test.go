package ztest

import (
	"fmt"
	"server/main/mmo_game/core"
	"testing"
)

//测试初始化的格子信息
func TestNewAOIManager(t *testing.T) {
	aoiMgr := core.NewAOIManager(100, 300, 4, 200, 450, 5)
	fmt.Println(aoiMgr)
}

//测试获取指定格子周边的九宫格信息
func TestAOIManagerSuroundGridsByGid(t *testing.T) {
	aoiMgr := core.NewAOIManager(0, 250, 5, 0, 250, 5)

	for k, _ := range aoiMgr.GetGrids() {
		//得到当前格子周边的九宫格
		grids := aoiMgr.GetSurroundGridsByGid(k)
		//得到九宫格所有的IDs
		fmt.Println("gid : ", k, " grids len = ", len(grids))
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Printf("grid ID: %d, surrounding grid IDs are %v\n", k, gIDs)
	}
}