import styles from "./stars.module.css";

const Starfield = () => {
  return (
    <div id="space">
      <div className={styles.stars}></div>
      <div className={styles.stars}></div>
      <div className={styles.stars}></div>
      <div className={styles.stars}></div>
      <div className={styles.stars}></div>
    </div>
  );
};

export default Starfield;
