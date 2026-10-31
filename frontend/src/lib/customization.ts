import { Protocols, type Protocol } from '@/schemas/primitives';

const allProtocols = Object.values(Protocols) as Protocol[];
const protocolSet = new Set<Protocol>(allProtocols);

export function panelBrandName(): string {
  return window.X_UI_BRAND_NAME?.trim() || '3X-UI';
}

export function panelBrandMark(): string {
  const words = panelBrandName().split(/\s+/u).filter(Boolean);
  if (words.length > 1) {
    return words
      .slice(0, 2)
      .map((word) => Array.from(word)[0])
      .join('')
      .toUpperCase();
  }
  return Array.from(words[0] || '3X')
    .slice(0, 2)
    .join('')
    .toUpperCase();
}

export function enabledProtocolValues(): Protocol[] {
  const configured = window.X_UI_ENABLED_PROTOCOLS;
  if (!configured) return allProtocols;

  const enabled = new Set<Protocol>();
  for (const value of configured) {
    const protocol = value.toLowerCase() as Protocol;
    if (protocolSet.has(protocol)) enabled.add(protocol);
  }
  return allProtocols.filter((protocol) => enabled.has(protocol));
}
