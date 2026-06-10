# 2S-UI Redesign — 设计系统使用手册

> redesign 分支:弃用 Vuetify,使用自有设计系统。本文档是所有页面重写的统一参照。
> 视觉真相:`../design_handoff_2s_ui_console/` 下的 JSX 文件(像素级对照);业务真相:旧 Vuetify 页面里的字段与 API 调用逻辑。

## 全局规则

- **禁止使用任何 `v-*` Vuetify 组件**。新页面只用本设计系统 + 原生元素。
- 颜色/间距/字号一律用 `src/styles/tokens.css` 的 CSS 变量(`var(--brand)`、`var(--surface-3)`、`var(--text-2)`、`var(--line)` …),**不得硬编码色值**(设计稿 JSX 中出现的 `#ff7eb6` 等特例除外,已收录在 `plugins/colors.ts`)。
- 等宽数字/代码一律加 `class="mono"`。
- 文案一律走 i18n:`$t('key')` / `useI18n().t`。优先复用现有 key(`actions.*`、`objects.*`、`client.*`、`in.*`、`out.*` 等,见 `src/locales/en.ts`);新视觉文案用 `ui.*` 前缀新 key(已在 locales 中补充)。
- 页面根元素:`<div class="page-stack fade-up">`(间距 16)或 `page-stack-lg`(18)。
- 业务数据/保存:沿用 Pinia `Data` store(`@/store/modules/data`):`loadData()`、`save(object, action, data, initUsers?)`;HTTP 用 `@/plugins/httputil`。**不要改 store 与 API 层。**

## 布局类(src/styles/app.css)

| class | 用途 |
|---|---|
| `.page-stack` / `.page-stack-lg` | 页面内容竖向堆叠(gap 16/18) |
| `.toolbar` | 页头操作条(flex、gap 10、wrap),`.grow` 占位 |
| `.entity-grid` | 卡片网格 `auto-fill minmax(248px,1fr)` gap 14 |
| `.grid2` | 表单两列 `1fr 1fr` gap 12 |
| `.field-grid` | 表单网格 `auto-fill minmax(180px,1fr)` gap 14 |
| `.form-divider` | 表单分隔线(`<hr class="form-divider">`) |
| `.dtable` | 数据表(thead 大写小标题、行 hover、`tr.clickable`) |
| `.kv-row` | 卡片内 key/value 行(`.k`/`.v`,`.v.mono`) |
| `.tabs` / `.tabs.page-tabs` | 下划线 tab(组件 `Tabs.vue` 已封装) |
| `.chip` + `.chip-emerald/rose/amber/brand/cyan` | 徽章(组件 `Chip.vue` 已封装) |
| `.input` | 文本框/选择框/textarea 通用(`<input class="input">`、`<select class="input">`、`<textarea class="input mono">`) |
| `.card` | 卡片容器 |

## UI 组件(src/components/ui/)

所有组件按需显式 import(无全局注册):`import Btn from '@/components/ui/Btn.vue'`

| 组件 | Props / 说明 |
|---|---|
| `Ico` | `name`(见 icons.ts 注册表)、`size=18`、`stroke=1.75` |
| `Btn` | `variant='primary'\|'ghost'(默认)\|'subtle'`、`sm`、`icon`(方形)、`loading`、`disabled`;slot 内容 `<Ico…/> 文本` |
| `Chip` | `color` = 预设名(emerald/rose/amber/brand/cyan)或任意 CSS 颜色(自动 tint);`dot` 加圆点 |
| `Toggle` | `v-model`(boolean)、`scale` |
| `Check` | `:checked` + `@toggle`(非 v-model,便于行内 stopPropagation) |
| `Segmented` | `v-model` + `options: [value,label][]`;`block` = drawer 内满宽样式 |
| `Tabs` | `v-model` + `tabs: [key,label][]`;`page` = 页级 tab(Settings) |
| `Field` | `label`、`hint`、`mb=15`;slot 放 `<input class="input">` 等 |
| `IconBtn` | 表格行内图标按钮:`name`、`danger`、`title` |
| `CardBtn` | 卡片底部动作按钮:`icon`、`label`(空=46px 方块)、`border`(左分隔线)、`danger` |
| `Avatar` | `name`、`size=30`(hash 渐变头像) |
| `Pop` | 弹出菜单容器。slot `#trigger="{ toggle, open }"` 放触发按钮,默认 slot `="{ close }"` 放 `.pop-item` 按钮列表;`width`、`align='end'`、`direction='down'` |
| `Drawer` | 右滑抽屉(RTL 自适应):`open`、`width=460`、`@close`;slots: `#title`、默认、`#footer` |
| `MDrawer` | 实体编辑抽屉(带图标头 + Cancel/Save 脚):`open`、`icon`、`color`、`title`、`sub`、`saveLabel`、`width=500`、`loading`、`@close`、`@save`;内容每次打开自动 remount |
| `Modal` | 居中模态:`open`、`title`、`width=720`、`@close`;slot 默认 + `#footer` |
| `DPanel` | 面板卡(头部 title/sub + #right 槽):`title`、`sub`、`pad=20`(0 = 内容无 padding,信息行列表用) |
| `EntityCard` | 实体卡:`title`、`type`、`color`、`icon`、`rows: {k,v,mono?,color?}[]`;slots `#chip`、`#actions`(放 CardBtn) |
| `RuleCard` | 规则卡:`n`、`logical`、`mode`、`action`、`targetKey`、`target`、`rules`、`invert`、`@edit`、`@del` |
| `SectionLabel` | 区块小标题(大写字距) |
| `MSwitchRow` | 圆角开关行(surface-3 背景):`v-model`、`label`、`desc` |
| `SwitchLabel` / `CheckLabel` | 行内 开关/复选 + 文字:`v-model`、`label` |
| `MultiPick` | chip 多选(selector/urltest 成员):`v-model: string[]`、`options: string[]` |
| `ChipSelect` | 下拉多选(chip 显示):`v-model: string[]`、`options:{title,value}[]`、`label` |
| `KeyInput` | 密钥输入 + 重新生成按钮:`v-model`、`placeholder`、`secret`、`@regenerate` |
| `MHint` | 虚线提示框(slot 文案) |
| `SettingsGroup` | Settings 卡片容器:`title?`、`desc?` |
| `SRow` | Settings 行(左 label 200px + 右控件 max 280px,移动端纵排):`label`、`hint` |
| `ToggleRow` | Settings 开关行:`v-model`、`label`、`desc` |
| `EmptyState` | 空状态:`icon`、`title`、`desc` |

### 共享业务组件(src/layouts/drawers/ 与 ui/)

| 组件 | 契约 |
|---|---|
| `layouts/drawers/StatsModal.vue` | 流量统计图模态:`:visible` `:resource`(inbound/outbound/client…对应 `objects.*` key)`:tag` + `@close`。内部走 `api/stats`,10s 自刷新 |
| `layouts/drawers/QrModal.vue` | 客户端二维码:`:visible` `:id`(client id)+ `@close` |
| `layouts/drawers/LogsModal.vue` / `BackupModal.vue` / `UsageStatsModal.vue` | `:visible` + `@close` |
| `layouts/drawers/EditorModal.vue` | 文本编辑模态:`:open` `:title` `:content` + `@save(value)` `@close` |
| `ui/DateTimeInput.vue` | 到期时间输入:`v-model`(**epoch 秒**,0=无限期);fa locale 自动用波斯历 |
| `plugins/clipboard.ts` | `copyToClipboard(txt, containerId?)`(带 toast) |

## 图表(src/components/charts/)

`Gauge`(径向仪表)、`AreaChart`(平滑面积图,可双系列 `data2` 虚线)、`Spark`(迷你折线)、`Donut`(`data:{value,color}[]`)、`BarMini`(进度条 `value` 0-100)。全部纯 SVG、theme-aware。

## 共享工具

- `@/plugins/colors`:`protoColor(type)`、`outColor(type)`、`dnsColor(type)`、`actionColor(action)`
- `@/plugins/utils`:`HumanReadable.sizeFormat(bytes)`(流量)、`formatSecond`、`remainedDays(epochSec)`
- `@/store/modules/app`:`theme`、`toggleTheme()`、`setLocale(l)`、`sidebarOpen`(移动端侧栏)

## Drawer 表单模式(所有编辑抽屉统一)

```vue
<MDrawer :open="open" icon="outbound" :color="outColor(form.type)"
  :title="isNew ? $t('ui.out.new') : item.tag" :sub="…"
  :save-label="isNew ? $t('actions.add') : $t('actions.save')"
  :loading="saving" @close="$emit('close')" @save="save">
  <Field :label="$t('type')">
    <select class="input" v-model="form.type">…</select>
  </Field>
  <div class="grid2"> <Field…/> <Field…/> </div>
</MDrawer>
```

- 打开时由父组件传入 `item`(null = 新建);MDrawer 自动 remount 内容,组件内 `watch(open)` 初始化本地 `form`(参照旧 Vuetify modal 的初始化逻辑)。
- 保存:组装 payload 后调 `Data().save(object, action, data)`,成功后 `$emit('close')`。**payload 结构必须与旧 modal 完全一致**(对照旧文件逐字段平移,含 `randomUtil` 的默认值生成)。

## 像素细节速查(来自设计稿)

- 卡片圆角 18(`--radius-lg`);按钮/输入框 10;chip 7。
- 主按钮 hover:`translateY(-1px)` + 加深辉光(tokens.css 已内置)。
- 抽屉宽:Client 460,实体编辑 500,规则 520;遮罩 `rgba(6,8,16,.5)`。
- 表格行高 ~42(th)/自适应(td padding `0 16px`,名称列 `11px 16px`)。
- 在线点:8px 圆 + `box-shadow: 0 0 0 4px color-mix(...18%)`。
- 协议色:vless=brand、vmess=violet、trojan=emerald、ss=amber、hy2=cyan、tuic=#ff7eb6、naive=#6ea8ff(用 `protoColor`)。
- 页面过渡 `.fade-up`;一切动效 `var(--ease)`、尊重 `prefers-reduced-motion`(tokens.css 已处理)。
