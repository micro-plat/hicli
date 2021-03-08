package tmpl

const MarkdownCurdSeqInstallGO = `package mysql

import "github.com/micro-plat/hydra"

func init() {
	hydra.Installer.DB.AddSQL(sys_sequence_info)
}

`

const MarkdownCurdSeqInstallSQL = `package mysql

const sys_sequence_info = {###}
DROP TABLE IF EXISTS sys_sequence_info;
CREATE TABLE sys_sequence_info  (
  seq_id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
  name varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '名称',
  create_time datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (seq_id) USING BTREE,
  INDEX idx_create_time(create_time) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '序列表' ROW_FORMAT = Compact;
{###}
`
