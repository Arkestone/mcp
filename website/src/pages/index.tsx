import type {ReactNode} from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';

import styles from './index.module.css';

const servers = [
  {name: 'mcp-instructions', desc: 'Copilot custom instruction files from local dirs and GitHub repos', port: '8080'},
  {name: 'mcp-skills', desc: 'Copilot skills with frontmatter metadata and reference bundles', port: '8081'},
  {name: 'mcp-prompts', desc: 'VS Code Copilot prompt and chat mode files', port: '8082'},
  {name: 'mcp-adr', desc: 'Architecture Decision Records from docs/adr/ directories', port: '8083'},
  {name: 'mcp-memory', desc: 'Persistent memory store for AI assistants across sessions', port: '8084'},
  {name: 'mcp-graph', desc: 'Knowledge graph with entity and relationship storage', port: '8085'},
];

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs">
            Get Started →
          </Link>
          <Link
            className="button button--outline button--secondary button--lg"
            href="https://github.com/Arkestone/mcp"
            style={{marginLeft: '1rem'}}>
            GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

function ServerCard({name, desc, port}: {name: string; desc: string; port: string}) {
  return (
    <div className="col col--4" style={{marginBottom: '1.5rem'}}>
      <div className="card shadow--md" style={{height: '100%'}}>
        <div className="card__header">
          <h3>
            <Link to={`/docs/servers/${name}`}>
              <code>{name}</code>
            </Link>
          </h3>
        </div>
        <div className="card__body">
          <p>{desc}</p>
        </div>
        <div className="card__footer">
          <small>Default port: <code>:{port}</code></small>
        </div>
      </div>
    </div>
  );
}

export default function Home(): ReactNode {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout title={siteConfig.title} description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <section style={{padding: '2rem 0'}}>
          <div className="container">
            <h2 style={{textAlign: 'center', marginBottom: '2rem'}}>Available Servers</h2>
            <div className="row">
              {servers.map((s) => (
                <ServerCard key={s.name} {...s} />
              ))}
            </div>
          </div>
        </section>
      </main>
    </Layout>
  );
}
