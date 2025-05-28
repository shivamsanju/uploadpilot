const cache: { [key: string]: string } = {};

const hashString = (str: string, length: number) => {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = (hash * 31 + str.charCodeAt(i)) % length;
  }
  return hash;
};

export const getRandomAvatarName = (key: string) => {
  const names = [
    'Destiny',
    'Jameson',
    'Andrea',
    'Luis',
    'Sara',
    'Amaya',
    'Emery',
    'Kimberly',
    'Liam',
    'Vivian',
    'Kingston',
    'Wyatt',
    'Eliza',
    'Robert',
    'Avery',
    'Sarah',
    'Maria',
    'Mason',
    'Oliver',
    'Leah',
  ];

  if (cache[key]) {
    return cache[key];
  }

  const index = hashString(key, names.length);

  const url = `https://api.dicebear.com/9.x/shapes/svg?seed=${names[index]}&backgroundColor=0a5b83,1c799f,69d2e7,f1f4dc,f88c49,b6e3f4,ffdfbf&backgroundType=gradientLinear`;

  cache[key] = url;

  return url;
};
