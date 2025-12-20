import { readFileSync, writeFileSync } from 'fs';
import { join } from 'path';

interface ChangelogEntry {
  version: string;
  date: string;
  changes: {
    added: string[];
    changed: string[];
    fixed: string[];
    removed: string[];
  };
}

function parseChangelog(content: string): ChangelogEntry[] {
  const entries: ChangelogEntry[] = [];
  const lines = content.split('\n');
  
  let currentEntry: ChangelogEntry | null = null;
  let currentSection: 'added' | 'changed' | 'fixed' | 'removed' | null = null;

  for (const line of lines) {
    // Match version header: ## [0.1.0] - 2024-01-15
    const versionMatch = line.match(/^##\s+\[([^\]]+)\]\s*-\s*(.+)$/);
    if (versionMatch) {
      if (currentEntry) {
        entries.push(currentEntry);
      }
      currentEntry = {
        version: versionMatch[1],
        date: versionMatch[2],
        changes: { added: [], changed: [], fixed: [], removed: [] }
      };
      currentSection = null;
      continue;
    }

    // Match section headers
    if (line.startsWith('### Added')) {
      currentSection = 'added';
      continue;
    }
    if (line.startsWith('### Changed')) {
      currentSection = 'changed';
      continue;
    }
    if (line.startsWith('### Fixed')) {
      currentSection = 'fixed';
      continue;
    }
    if (line.startsWith('### Removed')) {
      currentSection = 'removed';
      continue;
    }

    // Match bullet points
    if (currentEntry && currentSection && line.startsWith('- ')) {
      const item = line.substring(2).trim();
      if (item) {
        currentEntry.changes[currentSection].push(item);
      }
    }
  }

  // Push the last entry
  if (currentEntry) {
    entries.push(currentEntry);
  }

  return entries;
}

function generateAstroPage(entries: ChangelogEntry[]): string {
  const astroContent = `---
import Layout from '../../components/Layout.astro';
---

<Layout title="Changelog" description="Release history and changes for azd exec">
  <div class="page-container">
    <div class="page-header">
      <h1>Changelog</h1>
      <p class="page-intro">
        All notable changes to azd exec are documented here.
      </p>
      <p class="note">
        This page is auto-generated from 
        <a href="https://github.com/jongio/azd-exec/blob/main/cli/CHANGELOG.md" target="_blank" rel="noopener noreferrer">
          cli/CHANGELOG.md
        </a>
      </p>
    </div>

    <div class="content">
${entries.map(entry => generateEntryHTML(entry)).join('\n')}
    </div>
  </div>
</Layout>

<style>
  .page-container {
    max-width: 900px;
    margin: 0 auto;
    padding: 3rem 1.5rem;
  }

  .page-header {
    margin-bottom: 3rem;
  }

  .page-header h1 {
    font-size: 2.5rem;
    font-weight: 700;
    margin: 0 0 1rem;
    color: var(--color-text-primary);
  }

  .page-intro {
    font-size: 1.25rem;
    color: var(--color-text-secondary);
    margin: 0 0 0.5rem;
  }

  .note {
    font-size: 0.875rem;
    color: var(--color-text-tertiary);
    margin: 0;
  }

  .note a {
    color: var(--color-interactive-default);
    text-decoration: none;
  }

  .note a:hover {
    text-decoration: underline;
  }

  .content {
    line-height: 1.7;
  }

  .release {
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--color-border-default);
  }

  .release:last-child {
    border-bottom: none;
  }

  .release-header {
    display: flex;
    align-items: baseline;
    gap: 1rem;
    margin-bottom: 1.5rem;
  }

  .version {
    font-size: 1.75rem;
    font-weight: 600;
    color: var(--color-text-primary);
    font-family: var(--font-mono);
  }

  .date {
    font-size: 1rem;
    color: var(--color-text-tertiary);
  }

  .changes h3 {
    font-size: 1.125rem;
    font-weight: 600;
    margin: 1.5rem 0 0.75rem;
    color: var(--color-text-primary);
  }

  .changes ul {
    margin: 0 0 1rem;
    padding-left: 1.5rem;
  }

  .changes li {
    margin: 0.5rem 0;
    color: var(--color-text-secondary);
  }

  .badge {
    display: inline-block;
    font-size: 0.75rem;
    font-weight: 600;
    padding: 0.25rem 0.5rem;
    border-radius: 0.25rem;
    text-transform: uppercase;
    margin-right: 0.5rem;
  }

  .badge-added {
    background-color: #d1f4e0;
    color: #0d5c2f;
  }

  [data-theme="dark"] .badge-added {
    background-color: #0d5c2f;
    color: #d1f4e0;
  }

  .badge-changed {
    background-color: #d9e8ff;
    color: #0550ae;
  }

  [data-theme="dark"] .badge-changed {
    background-color: #0550ae;
    color: #d9e8ff;
  }

  .badge-fixed {
    background-color: #fff3cd;
    color: #664d03;
  }

  [data-theme="dark"] .badge-fixed {
    background-color: #664d03;
    color: #fff3cd;
  }

  .badge-removed {
    background-color: #ffebe9;
    color: #a0111f;
  }

  [data-theme="dark"] .badge-removed {
    background-color: #5a1a1a;
    color: #ffcccb;
  }
</style>
`;
  return astroContent;
}

function generateEntryHTML(entry: ChangelogEntry): string {
  const sections: string[] = [];

  if (entry.changes.added.length > 0) {
    sections.push(`        <h3><span class="badge badge-added">Added</span> New Features</h3>
        <ul>
${entry.changes.added.map(item => `          <li>${escapeHtml(item)}</li>`).join('\n')}
        </ul>`);
  }

  if (entry.changes.changed.length > 0) {
    sections.push(`        <h3><span class="badge badge-changed">Changed</span> Improvements</h3>
        <ul>
${entry.changes.changed.map(item => `          <li>${escapeHtml(item)}</li>`).join('\n')}
        </ul>`);
  }

  if (entry.changes.fixed.length > 0) {
    sections.push(`        <h3><span class="badge badge-fixed">Fixed</span> Bug Fixes</h3>
        <ul>
${entry.changes.fixed.map(item => `          <li>${escapeHtml(item)}</li>`).join('\n')}
        </ul>`);
  }

  if (entry.changes.removed.length > 0) {
    sections.push(`        <h3><span class="badge badge-removed">Removed</span> Deprecations</h3>
        <ul>
${entry.changes.removed.map(item => `          <li>${escapeHtml(item)}</li>`).join('\n')}
        </ul>`);
  }

  return `      <section class="release">
        <div class="release-header">
          <h2 class="version">${escapeHtml(entry.version)}</h2>
          <span class="date">${escapeHtml(entry.date)}</span>
        </div>
        <div class="changes">
${sections.join('\n')}
        </div>
      </section>`;
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

function main() {
  try {
    // Read changelog from CLI directory
    const changelogPath = join(process.cwd(), '..', 'cli', 'CHANGELOG.md');
    const changelogContent = readFileSync(changelogPath, 'utf-8');

    // Parse changelog
    const entries = parseChangelog(changelogContent);

    if (entries.length === 0) {
      console.warn('No changelog entries found');
      return;
    }

    // Generate Astro page
    const astroPage = generateAstroPage(entries);

    // Write to reference/changelog/index.astro
    const outputPath = join(process.cwd(), 'src', 'pages', 'reference', 'changelog', 'index.astro');
    writeFileSync(outputPath, astroPage, 'utf-8');

    console.log(`✅ Generated changelog page with ${entries.length} entries`);
    console.log(`   Output: ${outputPath}`);
  } catch (error) {
    console.error('❌ Error generating changelog:', error);
    process.exit(1);
  }
}

main();
