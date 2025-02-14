import React from "react";
import { Loader, Transition } from "@mantine/core";

export const FullScreenOverlay: React.FC<{ visible: boolean }> = ({
  visible,
}) => {
  return (
    <Transition mounted={visible} transition="fade" duration={300}>
      {(styles) => (
        <div style={styles}>
          <div
            style={{
              position: "fixed",
              top: 0,
              left: 0,
              width: "100%",
              height: "100%",
              backgroundColor:
                "light-dark(rgba(220, 220, 220, 0.2), rgba(19, 19, 19, 0.2))",
              backdropFilter: "blur(2px)",
              zIndex: 1000,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            <Loader />
          </div>
        </div>
      )}
    </Transition>
  );
};

export const ContainerOverlay: React.FC<{ visible: boolean }> = ({
  visible,
}) => {
  return (
    <Transition mounted={visible} transition="fade" duration={300}>
      {(styles) => (
        <div style={styles}>
          <div
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              width: "100%",
              height: "100%",
              backgroundColor:
                "light-dark(rgba(220, 220, 220, 0.2), rgba(19, 19, 19, 0.2))",
              backdropFilter: "blur(2px)",
              zIndex: 1000,
              display: "flex",
              justifyContent: "center",
              alignItems: "center",
              maskImage:
                "radial-gradient(circle, rgba(0, 0, 0, 0.5) 98%, rgba(0, 0, 0, 0) 100%)",
              WebkitMaskImage:
                "radial-gradient(circle, rgba(0, 0, 0, 0.5) 98%, rgba(0, 0, 0, 0) 100%)", // Ensures compatibility with WebKit browsers // Ensures compatibility with WebKit browsers
            }}
          >
            <Loader />
          </div>
        </div>
      )}
    </Transition>
  );
};
