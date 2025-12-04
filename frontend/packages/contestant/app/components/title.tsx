export function Title(props: { readonly children?: string }) {
  return (
    <title>
      {props.children != null ? `${props.children} | ICTSC2025` : "ICTSC2025"}
    </title>
  );
}
