# Feature Specification: 修复 TelegramSync 频道消息同步

**Feature Branch**: `006-fix-telegram-channel-sync`
**Created**: 2026-05-09
**Status**: Draft
**Input**: User description: "系统部署完成后，发现插件 TelegramSync 测试的时候只测试了群组，但是部署后，只绑定了两个频道，我没有发现产生文章"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 频道消息同步为文章 (Priority: P1)

系统管理员在后台配置了两个 Telegram 广播频道（broadcast channel），期望当频道发布新消息时，系统能自动将其同步为 CMS 文章。但部署后发现频道有新消息发布，系统中始终没有产生对应的文章。

**Why this priority**: 这是核心功能——频道消息同步是插件存在的主要目的，不工作意味着插件完全失效。

**Independent Test**: 可以通过在一个已绑定的 Telegram 广播频道中发布一条测试消息，然后检查 CMS 后台是否自动生成了对应的文章来验证。

**Acceptance Scenarios**:

1. **Given** 系统已绑定一个 Telegram 广播频道（broadcast channel）且状态为启用，**When** 该频道发布一条新文本消息，**Then** 系统应在 30 秒内自动创建对应的 CMS 文章，文章内容与消息一致
2. **Given** 系统已绑定一个 Telegram 广播频道，**When** 该频道发布包含图片的消息，**Then** 系统应创建文章并正确处理图片媒体
3. **Given** 系统同时绑定了群组和频道，**When** 群组和频道各有新消息，**Then** 两者的消息都应被正确同步为文章

---

### User Story 2 - 群组消息继续正常工作 (Priority: P2)

修复频道同步问题的同时，确保已有的群组（supergroup/megagroup）消息同步功能不受影响，继续正常工作。

**Why this priority**: 回归保护——不能因为修复频道而破坏已验证通过的群组功能。

**Independent Test**: 在已绑定的 Telegram 群组中发送一条消息，确认仍然能正常生成文章。

**Acceptance Scenarios**:

1. **Given** 系统已绑定一个 Telegram 群组（supergroup），**When** 群组中有成员发送新消息，**Then** 系统应正常创建对应的 CMS 文章
2. **Given** 系统同时绑定群组和频道，**When** 修复频道同步后，**Then** 群组消息同步不应受到任何影响

---

### User Story 3 - 频道配置不被意外清除 (Priority: P3)

频道的启用/禁用状态配置应该被正确维护，不应在处理消息时被意外修改或清除。

**Why this priority**: 数据完整性——配置丢失会导致已设置的频道监听失效，影响系统稳定性。

**Independent Test**: 配置多个频道（包括启用和禁用状态），触发消息处理后检查配置是否与预期一致。

**Acceptance Scenarios**:

1. **Given** 系统配置了 3 个频道（2 个启用、1 个禁用），**When** 任一启用频道收到消息并处理，**Then** 处理完成后 3 个频道的配置应保持不变（1 个禁用的频道仍然存在）
2. **Given** 一个频道被临时禁用后又重新启用，**When** 重新启用后该频道收到消息，**Then** 系统应正常同步消息

### Edge Cases

- 当 Telegram 广播频道被设置为私有频道时，消息同步是否正常？
- 当频道消息为纯转发消息（forward）时，应正常同步为文章（与普通消息一致）
- 当客户端与 Telegram 服务器断开后重连，系统应自动恢复断线期间遗漏的频道消息
- 当频道发送空消息（仅含媒体无文本）时，文章应如何创建？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 插件必须能正确接收并处理来自 Telegram 广播频道（broadcast channel）的新消息，将其同步为 CMS 文章
- **FR-002**: 插件必须能正确接收并处理来自 Telegram 群组（supergroup/megagroup）的新消息，维持现有功能不受影响
- **FR-003**: 系统必须正确区分频道消息和群组消息的来源，使用正确的消息路由路径
- **FR-004**: 消息处理过程中不得修改或丢失频道配置数据（包括已禁用的频道记录）
- **FR-005**: 当 Telegram 客户端通过 `updates.Manager` 接收更新时，必须确保 `UpdateNewChannelMessage` 类型的更新被正确分发到对应的处理函数
- **FR-006**: 当消息来自广播频道时，频道 ID 必须与后台配置中的频道 ID 正确匹配，不因 ID 格式差异导致匹配失败
- **FR-007**: 插件应在消息处理失败时记录明确的日志信息，包括更新类型、频道 ID、失败原因，便于排查问题
- **FR-008**: 当 Telegram 客户端断线后重连，系统应自动尝试恢复断线期间遗漏的频道消息，避免内容丢失

### Key Entities

- **TelegramChannel（频道配置）**: 记录已绑定的 Telegram 频道/群组信息，包含频道 ID、名称、状态（启用/禁用）、同步规则、统计信息等
- **TelegramSyncLog（同步日志）**: 记录每次消息同步的结果，包含频道 ID、消息 ID、文章 ID、状态（成功/失败/跳过）和错误信息
- **ChannelConfig（运行时配置）**: 从 JSON 解析的频道配置，用于消息处理时查找匹配的频道

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 在 Telegram 广播频道发布消息后，系统在 60 秒内自动创建对应的 CMS 文章
- **SC-002**: 群组消息同步功能保持 100% 不受影响，与修复前行为一致
- **SC-003**: 消息处理后，所有频道配置（包括禁用状态的频道）保持不变
- **SC-004**: 系统日志中能清晰区分频道消息和群组消息的处理路径，便于问题排查
- **SC-005**: 连续运行 24 小时后，频道消息同步不出现遗漏或重复

## Clarifications

### Session 2026-05-09

- Q: 频道转发消息应如何处理？ → A: 正常同步为文章（与普通消息一致）
- Q: 断线重连后的消息恢复策略？ → A: 自动恢复：重连后尝试补回断线期间遗漏的消息

## Assumptions

- Telegram 广播频道和群组在 gotd/td 库中使用不同的更新类型（`UpdateNewChannelMessage` vs `UpdateNewMessage`）
- 生产环境使用带 session storage 的客户端路径（`NewClientWithStorage`），即通过 `updates.Manager` 和 `UpdateDispatcher` 分发消息
- 频道通过后台管理界面配置，频道 ID 来源于 Telegram API 返回的 `channel.ID` 字段
- 已有的群组消息同步功能经过测试验证是正常的，只需确保频道消息走正确的路径
