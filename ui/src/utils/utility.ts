export const formatBytes = (bytes: number, decimals: number = 2) => {
  if (!+bytes) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = [
    "Bytes",
    "KiB",
    "MiB",
    "GiB",
    "TiB",
    "PiB",
    "EiB",
    "ZiB",
    "YiB",
  ];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
};

export const areBracketsBalanced = (input: string) => {
  const stack = [];
  const brackets: { [key: string]: string } = {
    ")": "(",
    "}": "{",
    "]": "[",
  };

  for (const char of input) {
    if (["(", "{", "["].includes(char)) {
      stack.push(char);
    } else if ([")", "}", "]"].includes(char)) {
      if (stack.pop() !== brackets[char]) {
        return false;
      }
    }
  }

  return stack.length === 0;
};

export const getFilterParams = (filter: Record<string, string[]>): string => {
  let filterParam = "";
  if (Object.entries(filter).length > 0) {
    const fp: string[] = [];
    for (const [key, values] of Object.entries(filter)) {
      const encodedValues = values.map((value) => encodeURIComponent(value));
      if (values.length !== 0) {
        fp.push(`${key}:${encodedValues.join(",")}`);
      }
    }

    if (fp.length !== 0) {
      filterParam = `&filter=${fp.join(";")}`;
    }
  }

  return filterParam;
};

export const getSearchParam = (search: string): string => {
  return search ? `&search=${encodeURIComponent(search)}` : "";
};
