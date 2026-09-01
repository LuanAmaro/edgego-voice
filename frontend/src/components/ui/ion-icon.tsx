import React from "react";

interface IonIconProps extends React.HTMLAttributes<HTMLElement> {
  name: string;
  className?: string;
  style?: React.CSSProperties;
}

export function IonIcon({ name, className = "", style, ...props }: IonIconProps) {
  // @ts-expect-error ion-icon is a web component
  return <ion-icon name={name} class={className} style={style} {...props} />;
}
