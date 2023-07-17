const verifyOutlook = async (address: string): Promise<string> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return "";
};

const verifyIMAP = async (
  address: string,
  password: string
): Promise<string> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return "";
};
const verifyPOP3 = async (
  address: string,
  password: string
): Promise<string> => {
  await new Promise((resolve) => setTimeout(resolve, 1000));
  return "";
};

export { verifyOutlook, verifyIMAP, verifyPOP3 };
