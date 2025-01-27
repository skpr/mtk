import clsx from 'clsx';
import src from './how-it-works.png';

export default function HowItWorks() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className={clsx('col col--8 col--offset-2')}>
            <img src={src} alt="Overview of how MTK works" />;
          </div>
        </div>
      </div>
    </section>
  );
}
