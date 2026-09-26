# Domainry Knowledge SDK

Public Knowledge contracts and in-process module boundaries.

This module contains no database implementation. A product composition root
selects a Knowledge implementation through `modulehost.Factory`; the
implementation creates and owns its Store behind `modulehost.ModuleBinding`.

`files` 是独立于 Agent 会话模型的通用文件契约和远程客户端。产品通过
`files.OpenRemote` 连接独立 Knowledge 服务，只保存返回的 `file_id`；文件上传不会自动
进入知识库或触发索引。
