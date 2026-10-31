import { afterEach, expect, test } from 'vitest';

import { enabledProtocolValues, panelBrandMark, panelBrandName } from '@/lib/customization';
import { Protocols } from '@/schemas/primitives';

afterEach(() => {
  delete window.X_UI_BRAND_NAME;
  delete window.X_UI_ENABLED_PROTOCOLS;
});

test('uses the upstream brand and all protocols when no private profile is configured', () => {
  expect(panelBrandName()).toBe('3X-UI');
  expect(enabledProtocolValues()).toEqual(Object.values(Protocols));
});

test('normalizes the private brand and protocol allowlist from runtime globals', () => {
  window.X_UI_BRAND_NAME = '  Private Proxy  ';
  window.X_UI_ENABLED_PROTOCOLS = ['vless', 'hysteria', 'tuic', 'vless', 'unknown'];

  expect(panelBrandName()).toBe('Private Proxy');
  expect(panelBrandMark()).toBe('PP');
  expect(enabledProtocolValues()).toEqual(['vless', 'hysteria', 'tuic']);
});
