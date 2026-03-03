import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Servers',
      collapsed: false,
      items: [
        'servers/mcp-instructions',
        'servers/mcp-skills',
        'servers/mcp-prompts',
        'servers/mcp-adr',
        'servers/mcp-memory',
        'servers/mcp-graph',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'network',
        'contributing',
        'security',
        'changelog',
      ],
    },
  ],
};

export default sidebars;
