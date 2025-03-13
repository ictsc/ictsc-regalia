export function Title(props: { readonly children?: string }) {
  return (
    <title>
      {props.children != null ? `${props.children} | ICTSC2024` : "ICTSC2024"}
    </title>
  );
}
