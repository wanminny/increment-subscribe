package canal

import (
	"sync"

	"github.com/siddontang/go-mysql/mysql"
	"gopkg.in/birkirb/loggers.v1/log"

	"fmt"
	. "lt-test/supplier/env"
	"lt-test/supplier/tools"
)

type masterInfo struct {
	sync.RWMutex

	pos mysql.Position

	gtid mysql.GTIDSet
}

func (m *masterInfo) Update(pos mysql.Position) {
	log.Debugf("update master position %s", pos)

	m.Lock()
	m.pos = pos
	m.Unlock()
}

func (m *masterInfo) UpdateGTID(gtid mysql.GTIDSet) {
	log.Debugf("update master gtid %s", gtid.String())

	m.Lock()
	m.gtid = gtid
	//记录文件
	binInfo := fmt.Sprintf("%s,%s\n", tools.CurrentTime(), m.gtid)
	tools.SaveToFile(binInfo, BIN_LOG_FILE_TO_READ_GTID)
	m.Unlock()
}

func (m *masterInfo) Position() mysql.Position {
	m.RLock()
	defer m.RUnlock()

	return m.pos
}

func (m *masterInfo) GTID() mysql.GTIDSet {
	m.RLock()
	defer m.RUnlock()

	return m.gtid
}
