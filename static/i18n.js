const I18N_MAP = {
  'Logs': {
    'en': 'Logs',
    'zh-cn': '日志'
  },
  'Save': {
    'en': 'Save',
    'zh-cn': '保存'
  },
  'Config:': {
    'en': 'Config:',
    'zh-cn': '配置切换:'
  },
  'Add': {
    'en': 'Add',
    'zh-cn': '添加'
  },
  'Rename': {
    'en': 'Rename',
    'zh-cn': '重命名'
  },
  'RenameHelp': {
    'en': 'Enter a new name:',
    'zh-cn': '输入新名称：'
  },
  'Delete': {
    'en': 'Delete',
    'zh-cn': '删除'
  },
  'DNS Provider': {
    'en': 'DNS Provider',
    'zh-cn': 'DNS服务商'
  },
  'Create AccessKey': {
    'en': 'Create AccessKey',
    'zh-cn': '创建 AccessKey'
  },
  'Auto': {
    'en': 'Auto',
    'zh-cn': '自动'
  },
  '1s': {
    'en': '1s',
    'zh-cn': '1秒'
  },
  '5s': {
    'en': '5s',
    'zh-cn': '5秒'
  },
  '10s': {
    'en': '10s',
    'zh-cn': '10秒'
  },
  '1m': {
    'en': '1m',
    'zh-cn': '1分钟'
  },
  '2m': {
    'en': '2m',
    'zh-cn': '2分钟'
  },
  '10m': {
    'en': '10m',
    'zh-cn': '10分钟'
  },
  '30m': {
    'en': '30m',
    'zh-cn': '30分钟'
  },
  '1h': {
    'en': '1h',
    'zh-cn': '1小时'
  },
  'ttlHelp': {
    'en': 'You can modify it if the account supports a smaller TTL. The TTL will only be updated when the IP changes',
    'zh-cn': '如账号支持更小的 TTL, 可修改。IP 有变化时才会更新TTL'
  },
  'Enabled': {
    'en': 'Enabled',
    'zh-cn': '是否启用'
  },
  'Get IP method': {
    'en': 'Get IP method',
    'zh-cn': '获取 IP 方式'
  },
  'By api': {
    'en': 'By api',
    'zh-cn': '通过接口获取'
  },
  'By network card': {
    'en': 'By network card',
    'zh-cn': '通过网卡获取'
  },
  'By command': {
    'en': 'By command',
    'zh-cn': '通过命令获取'
  },
  'Expand': {
    'en': 'Expand',
    'zh-cn': '展开'
  },
  'Collapse': {
    'en': 'Collapse',
    'zh-cn': '收起'
  },
  'domainsHelp': {
    'en': `
      Enter one domain per line.
      If the domain is unregistrable, manually separate it into a subdomain and a root domain by using a colon. e.g. <code>www:domain.example.com</code><br />

      Support for <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">custom parameters</a> (Simplified Chinese)
    `,
    'zh-cn': `
      每行一个域名。
      如果域名不可注册，请使用冒号手动将其分为子域名和根域名。如 <code>www:domain.example.com</code><br />
      支持<a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/传递自定义参数">自定义参数</a>
    `
  },
  'Regular exp.': {
    'en': 'Regular exp.',
    'zh-cn': '匹配正则表达式'
  },
  'regHelp': {
    'en': 'You can use @1 to specify the first IPv6 address, @2 to specify the second IPv6 address... You can also use regular expressions to match the specified IPv6 address, leave it blank to disable it',
    'zh-cn': '可使用 @1 指定第一个IPv6地址, @2 指定第二个IPv6地址... 也可使用正则表达式匹配指定的IPv6地址, 留空则不启用'
  },
  'Others': {
    'en': 'Others',
    'zh-cn': '其他'
  },
  'Deny from WAN': {
    'en': 'Deny from WAN',
    'zh-cn': '禁止公网访问'
  },
  'NotAllowWanAccessHelp': {
    'en': 'Enable to deny access from the public network',
    'zh-cn': '启用后禁止从公网访问此页面'
  },
  'Username': {
    'en': 'Username',
    'zh-cn': '用户名'
  },
  'accountHelp': {
    'en': 'Username/Password is required',
    'zh-cn': '必须输入用户名/密码'
  },
  'passwordHelp': {
    'en': 'If you need to change the password, please enter it here',
    'zh-cn': '如需修改密码，请在此处输入新密码'
  },
  'Password': {
    'en': 'Password',
    'zh-cn': '密码'
  },
  'WebhookURLHelp': {
    'en': `
      <a
        target="blank"
        href="https://github.com/jeessy2/ddns-go/blob/master/README_EN.md#webhook"
      >Click to get more info</a
      ><br />
      Support variables #{ipv4Addr}, #{ipv4Result},
      #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}, #{timestamp}
    `,
    'zh-cn': `
      <a target="blank" href="https://github.com/jeessy2/ddns-go#webhook">点击参考官方 Webhook 说明</a>
      <br />
      支持的变量 #{ipv4Addr}, #{ipv4Result}, #{ipv4Domains}, #{ipv6Addr}, #{ipv6Result}, #{ipv6Domains}, #{timestamp}
    `
  },
  'WebhookRequestBodyHelp': {
    'en': 'Leave empty for a GET request; if filled in, a POST request is used. Supported placeholders: <code>#{ipv4Addr}</code> <code>#{ipv6Addr}</code> <code>#{oldIpv6Addr}</code> <code>#{timestamp}</code>',
    'zh-cn': '留空为 GET 请求，填写则使用 POST。支持占位符： <code>#{ipv4Addr}</code> <code>#{ipv6Addr}</code> <code>#{oldIpv6Addr}</code> <code>#{timestamp}</code>'
  },
  'WebhookHeadersHelp': {
    'en': 'One header per line, such as: Authorization: Bearer API_KEY',
    'zh-cn': '一行一个Header, 如: Authorization: Bearer API_KEY'
  },
  'Try it': {
    'en': 'Try it',
    'zh-cn': '模拟测试Webhook'
  },
  'Clear': {
    'en': 'Clear',
    'zh-cn': '清空'
  },
  'OK': {
    'en': 'OK',
    'zh-cn': '确定'
  },
  'languageSelectTooltip': {
    'en': 'Select language',
    'zh-cn': '选择语言'
  },
  "Ipv4UrlHelp": {
    'en': "https://api.ipify.org, https://ddns.oray.com/checkip, https://ip.3322.net, https://v4.yinghualuo.cn/bejson, https://myip.ipip.net",
    'zh-cn': "https://ddns.oray.com/checkip, https://ip.3322.net, https://v4.yinghualuo.cn/bejson, https://myip.ipip.net"
  },
  "Ipv6UrlHelp": {
    'en': "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn, https://v6.yinghualuo.cn/bejson",
    'zh-cn': "https://speed.neu6.edu.cn/getIP.php, https://v6.ident.me, https://6.ipw.cn, https://v6.yinghualuo.cn/bejson"
  },
  "Ipv4NetInterfaceHelp": {
    'en': "Get IPv4 address through network card",
    'zh-cn': "通过网卡获取IPv4"
  },
  "Ipv6NetInterfaceHelp": {
    'en': "If you do not specify a matching regular expression, the first IPv6 address will be used by default",
    'zh-cn': "如不指定匹配正则表达式，将默认使用第一个 IPv6 地址"
  },
  "Ipv4CmdHelp": {
    'en': "Get IPv4 through command, only use the first matching IPv4 address of standard output(stdout). Such as: ip -4 addr show eth1",
    'zh-cn': `
      通过命令获取IPv4, 仅使用标准输出(stdout)的第一个匹配的 IPv4 地址。如: ip -4 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `
  },
  "Ipv6CmdHelp": {
    'en': "Get IPv6 through command, only use the first matching IPv6 address of standard output(stdout). Such as: ip -6 addr show eth1",
    'zh-cn': `
      通过命令获取IPv6, 仅使用标准输出(stdout)的第一个匹配的 IPv6 地址。如: ip -6 addr show eth1
      <a target="blank" href="https://github.com/jeessy2/ddns-go/wiki/通过命令获取IP参考">点击参考更多</a>
    `
  },
  "NetInterfaceEmptyHelp": {
    'en': '<span style="color: red">No available network card found</span>',
    'zh-cn': '<span style="color: red">没有找到可用的网卡</span>'
  },
  "Http Interface": {
    'en': 'Http Interface',
    'zh-cn': 'HTTP 请求网卡'
  },
  "Default": {
    'en': 'Default',
    'zh-cn': '默认'
  },
  "HttpInterfaceHelp": {
    'en': 'Bind HTTP requests to a specific network interface (similar to curl --interface). Leave empty to use the default.',
    'zh-cn': '发送 HTTP 请求时绑定指定网卡（类似 curl --interface）。留空则使用默认网卡。'
  },
  "Login": {
    'en': 'Login',
    'zh-cn': '登录'
  },
  "LoginInit": {
    'en': 'Login and configure as an administrator account',
    'zh-cn': '登录并配置为管理员账号'
  },
  "Logout": {
    'en': 'Logout',
    'zh-cn': '注销'
  },
  "webhookTestTooltip": {
    'en': 'Send a fake data to the Webhook URL immediately to test if the Webhook is working properly',
    'zh-cn': '立即发送一条假数据到Webhook URL，用于测试Webhook是否正常工作'
  },
  "themeTooltip": {
    'en': 'Click: Switch theme<br>Long press: Restore auto mode',
    'zh-cn': '单击：切换明暗主题<br>长按：恢复自动跟随系统'
  },
  "extParamHelp": {
    'en': 'Optional. If you are using a Vercel Team account, please fill in the Team ID',
    'zh-cn': '可选项，如果您使用的是 Vercel 团队账户，请填写团队 ID'
  },
  // writing.html - 运行状态
  'StatusTitle': {
    'en': 'Status',
    'zh-cn': '运行状态'
  },
  'StatusCurrentIPv4': {
    'en': 'Current IPv4',
    'zh-cn': '当前 IPv4'
  },
  'StatusCurrentIPv6': {
    'en': 'Current IPv6',
    'zh-cn': '当前 IPv6'
  },
  'StatusLastCheck': {
    'en': 'Last check',
    'zh-cn': '上次检测'
  },
  'StatusLastIPv6Change': {
    'en': 'Last IPv6 change',
    'zh-cn': 'IPv6 最近变化'
  },
  'StatusCertIPs': {
    'en': 'Certificate IPs',
    'zh-cn': '证书包含地址'
  },
  'StatusCertNotAfter': {
    'en': 'Certificate expiry',
    'zh-cn': '证书到期时间'
  },
  'StatusCertTask': {
    'en': 'Certificate task',
    'zh-cn': '证书任务'
  },
  'StatusNextRenew': {
    'en': 'Next auto renewal',
    'zh-cn': '下次自动更新'
  },
  'StatusLastResult': {
    'en': 'Last result',
    'zh-cn': '最近结果'
  },
  'CheckNowBtn': {
    'en': 'Check now',
    'zh-cn': '立即检测'
  },
  'RenewCertBtn': {
    'en': 'Renew certificate now',
    'zh-cn': '立即更新证书'
  },
  'ManualRenewHelp': {
    'en': "Manual renewal forces certificate renewal. Please be aware of Let's Encrypt rate limits.",
    'zh-cn': "手动更新会强制续期证书，请注意 Let's Encrypt 的频率限制。"
  },
  // writing.html - IP 检测
  'IpCheckTitle': {
    'en': 'IP Detection',
    'zh-cn': 'IP 检测'
  },
  'IntervalLabel': {
    'en': 'Check interval',
    'zh-cn': '检测频率'
  },
  'SecondsUnit': {
    'en': 's',
    'zh-cn': '秒'
  },
  'IntervalHelp': {
    'en': 'Minimum 5 seconds. The new interval takes effect immediately after saving.',
    'zh-cn': '最小 5 秒，保存后立即以新频率检测。'
  },
  'EnableLabel': {
    'en': 'Enable',
    'zh-cn': '启用'
  },
  'GetTypeLabel': {
    'en': 'Get method',
    'zh-cn': '获取方式'
  },
  'NetInterfaceOption': {
    'en': 'Network interface',
    'zh-cn': '网卡'
  },
  'UrlOption': {
    'en': 'URL',
    'zh-cn': '接口'
  },
  'CmdOption': {
    'en': 'Command',
    'zh-cn': '命令'
  },
  'UrlPlaceholder': {
    'en': 'Interface URLs, separated by commas',
    'zh-cn': '接口地址，多个用逗号分隔'
  },
  'Ipv4CmdPlaceholder': {
    'en': 'Command to get IPv4',
    'zh-cn': '获取 IPv4 的命令'
  },
  'Ipv6CmdPlaceholder': {
    'en': 'Command to get IPv6',
    'zh-cn': '获取 IPv6 的命令'
  },
  'Ipv6MatchLabel': {
    'en': 'IPv6 match',
    'zh-cn': 'IPv6 匹配'
  },
  'Ipv6MatchPlaceholder': {
    'en': 'Leave empty to use the first address; @2 means the second address; a regular expression is also supported',
    'zh-cn': '留空使用第一个地址；@2 表示第二个地址；也可填正则表达式'
  },
  'Ipv6MatchHelp': {
    'en': 'When the network interface has multiple IPv6 addresses, use @index or a regular expression to select one.',
    'zh-cn': '网卡有多个 IPv6 时，可用 @序号 或正则表达式选择。'
  },
  // writing.html - Webhook
  'WebhookTitle': {
    'en': 'Webhook Notification',
    'zh-cn': 'Webhook 通知'
  },
  'WebhookUrlLabel': {
    'en': 'URL',
    'zh-cn': 'URL'
  },
  'WebhookUrlPlaceholder': {
    'en': 'URL to request when IPv6 changes',
    'zh-cn': 'IPv6 变化时请求的地址'
  },
  'WebhookRequestBodyLabel': {
    'en': 'RequestBody',
    'zh-cn': 'RequestBody'
  },
  'WebhookRequestBodyPlaceholder': {
    'en': 'e.g. {"ipv6":"#{ipv6Addr}","old":"#{oldIpv6Addr}"}',
    'zh-cn': '例如: {"ipv6":"#{ipv6Addr}","old":"#{oldIpv6Addr}"}'
  },
  'WebhookHeadersLabel': {
    'en': 'Headers',
    'zh-cn': 'Headers'
  },
  'WebhookHeadersPlaceholder': {
    'en': 'One per line, e.g. Authorization: Bearer xxx',
    'zh-cn': '每行一个，例如: Authorization: Bearer xxx'
  },
  'TestWebhookBtn': {
    'en': 'Test Webhook',
    'zh-cn': '测试 Webhook'
  },
  'TestWebhookHelp': {
    'en': 'Send one request with fake data. Please save the actual configuration first.',
    'zh-cn': '使用假数据发送一次请求，实际配置请先保存。'
  },
  // writing.html - Certbot 证书
  'CertTitle': {
    'en': 'Certbot Certificate',
    'zh-cn': 'Certbot 证书'
  },
  'CertAutoIssueLabel': {
    'en': 'Auto issue',
    'zh-cn': '自动申请'
  },
  'CertAutoIssueHelp': {
    'en': 'Automatically request a certificate when IPv6 changes or the certificate does not match, and restart Nginx.',
    'zh-cn': 'IPv6 变化或证书不匹配时自动申请证书，并重启 Nginx。'
  },
  'CertbotPathLabel': {
    'en': 'certbot path',
    'zh-cn': 'certbot 路径'
  },
  'CertbotPathPlaceholder': {
    'en': 'certbot or /usr/bin/certbot',
    'zh-cn': 'certbot 或 /usr/bin/certbot'
  },
  'CertEmailLabel': {
    'en': 'Email',
    'zh-cn': '邮箱'
  },
  'CertEmailPlaceholder': {
    'en': 'Optional, used for certificate expiry reminders',
    'zh-cn': '可选，用于证书到期提醒'
  },
  'CertChallengeLabel': {
    'en': 'Challenge type',
    'zh-cn': '验证方式'
  },
  'CertChallengeStandalone': {
    'en': 'standalone (temporarily listen on port 80)',
    'zh-cn': 'standalone（临时监听 80 端口）'
  },
  'CertChallengeWebroot': {
    'en': 'webroot (use an existing website directory)',
    'zh-cn': 'webroot（使用已有网站目录）'
  },
  'CertWebrootLabel': {
    'en': 'Webroot',
    'zh-cn': 'Webroot'
  },
  'CertProfileLabel': {
    'en': 'Certificate Profile',
    'zh-cn': 'Certificate Profile'
  },
  'CertProfileHelp': {
    'en': "Let's Encrypt IP certificates require <code>shortlived</code> (certbot 5.3+).",
    'zh-cn': "Let's Encrypt 的 IP 证书需使用 <code>shortlived</code>（certbot 5.3+）。"
  },
  'CertLinkDirLabel': {
    'en': 'Certificate directory',
    'zh-cn': '证书目录'
  },
  'CertLinkDirHelp': {
    'en': 'The Nginx certificate for port 8443 should point to <code>directory/fullchain.pem</code> and <code>directory/privkey.pem</code>. After issuance or renewal, the certificate is automatically copied to and overwrites the files here (the directory is created if it does not exist).',
    'zh-cn': 'Nginx 8443 端口证书应指向 <code>该目录/fullchain.pem</code> 与 <code>该目录/privkey.pem</code>， 证书申请/续期后会自动复制并覆盖此处文件（目录不存在会自动创建）。'
  },
  'CertOwnerLabel': {
    'en': 'Certificate owner user',
    'zh-cn': '证书所属用户'
  },
  'CertOwnerPlaceholder': {
    'en': 'e.g. pengge (leave empty so that only root can read the private key)',
    'zh-cn': '例如: pengge（留空则仅 root 可读私钥）'
  },
  'CertOwnerHelp': {
    'en': 'Fill in the user that runs Nginx when it runs as non-root (e.g. <code>brew services start nginx</code> on macOS). The program will chown the certificate to that user as root to avoid permission errors when loading the certificate.',
    'zh-cn': 'Nginx 以非 root 运行（如 macOS 的 <code>brew services start nginx</code>）时填写其运行用户， 程序会以 root 身份把证书 chown 给该用户，避免加载证书时权限不足。'
  },
  'CertNginxReloadLabel': {
    'en': 'Restart Nginx command',
    'zh-cn': '重启 Nginx 命令'
  },
  'CertExtraArgsLabel': {
    'en': 'Extra arguments',
    'zh-cn': '附加参数'
  },
  'CertExtraArgsPlaceholder': {
    'en': 'e.g. --key-type ecdsa',
    'zh-cn': '例如: --key-type ecdsa'
  },
  'CertDailyRenewLabel': {
    'en': 'Daily scheduled renewal',
    'zh-cn': '每日定时更新'
  },
  'CertRenewTimeLabel': {
    'en': 'Time',
    'zh-cn': '时间'
  },
  'HourUnit': {
    'en': 'h',
    'zh-cn': '时'
  },
  'MinuteUnit': {
    'en': 'min',
    'zh-cn': '分'
  },
  'CertDeleteOldLabel': {
    'en': 'Delete old certificates',
    'zh-cn': '删除旧证书'
  },
  'CertDeleteOldHelp': {
    'en': 'Delete the certificates for old IPv6 addresses after the new certificate is deployed successfully.',
    'zh-cn': '新证书部署成功后删除旧 IPv6 地址对应的证书。'
  },
  // writing.html - 其它
  'OthersTitle': {
    'en': 'Others',
    'zh-cn': '其它'
  },
  'PasswordPlaceholder': {
    'en': 'Leave empty to keep the current password',
    'zh-cn': '留空则不修改密码'
  },
};

const getCurrentLang = () => {
  const serverLang = window.__DDNS_GO_LANG__;
  if (serverLang) {
    return String(serverLang).replaceAll('_', '-').toLowerCase();
  }
  return localStorage.getItem('lang') || (navigator.language || navigator.browserLanguage).replaceAll('_', '-').toLowerCase();
};

const getLocalLang = (langs) => {
  const langValue = getCurrentLang();
  if (langs.includes(langValue)) {
    return langValue;
  }
  if (langs.includes(langValue.split('-')[0])) {
    return langValue.split('-')[0];
  }
  for (const l of langs) {
    if (l.split('-')[0] === langValue.split('-')[0]) {
      return l;
    }
  }
  return 'en';
}

const i18n = (keyOrLangDict) => {
  let key = keyOrLangDict;
  let langDict = keyOrLangDict;
  if (typeof keyOrLangDict === 'string') {
    langDict = I18N_MAP[keyOrLangDict];
  } else {
    key = null;
  }
  if (!langDict) {
    console.warn(`i18n: No translation for key "${key}"`);
    return key;
  }
  const lang = getLocalLang(Object.keys(langDict));
  if (lang in langDict) {
    return langDict[lang];
  }
  console.warn(`i18n: No such language "${lang}" in langDict ${langDict}`);
  return key;
}

const convertDom = (dom = document) => {
  dom.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.dataset.i18n;
    el.textContent = i18n(key);
  });
  dom.querySelectorAll('[data-i18n-html]').forEach(el => {
    const key = el.dataset.i18nHtml;
    el.innerHTML = i18n(key);
  });
  dom.querySelectorAll('[data-i18n-attr]').forEach(el => {
    el.dataset.i18nAttr.split(',').forEach(item => {
      let [attr, key] = item.split(':');
      attr = attr.trim();
      key = key || el.getAttribute(attr);
      el.setAttribute(attr, i18n(key));
    });
  });
}

document.addEventListener('DOMContentLoaded', () => { convertDom(); });
