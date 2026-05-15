// Минимальный клиент Connect-RPC поверх fetch.
//
// Connect поддерживает JSON-кодирование поверх POST: достаточно отправить
// тело как обычный JSON и принять JSON в ответ. Генерация TS-клиента не
// нужна — типизация описана локально и совпадает с .proto.

export type GetCompanyResponse = {
  company?: Company;
};

export type Company = {
  securityDescription?: Record<string, unknown>;
  stock?: Stock;
};

export type Stock = {
  info?: Record<string, unknown>;
  summary?: Record<string, unknown>;
  ratios?: Record<string, unknown>[];
  reports?: Record<string, unknown>[];
  dividends?: Record<string, unknown>[];
  ideas?: Record<string, unknown>[];
  insiderTransactions?: Record<string, unknown>[];
  operations?: Record<string, unknown>[];
  owners?: Record<string, unknown>[];
  shares?: Record<string, unknown>[];
};

export type ConnectError = {
  code: string;
  message: string;
};

export async function getCompany(ticker: string): Promise<GetCompanyResponse> {
  const resp = await fetch('/company.v1.CompanyService/GetCompany', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
    body: JSON.stringify({ ticker }),
  });

  if (!resp.ok) {
    let payload: ConnectError | null = null;
    try {
      payload = (await resp.json()) as ConnectError;
    } catch {
      // ignore
    }
    const code = payload?.code ?? `http_${resp.status}`;
    const message = payload?.message ?? resp.statusText;
    throw new Error(`${code}: ${message}`);
  }

  return (await resp.json()) as GetCompanyResponse;
}
