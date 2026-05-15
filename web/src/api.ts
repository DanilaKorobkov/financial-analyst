// Connect-RPC клиент поверх типов, сгенерированных из proto-схем
// (артефакты `task web:generate-proto`, лежат в src/gen/). Ручного описания
// контракта здесь нет — расхождение схемы и клиента ловится в CI задачей
// `web:lint-generated-proto`.
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { CompanyService } from './gen/company/v1/company_pb';

// baseUrl относительный: в dev Vite перенаправляет /company.v1.* на backend
// (см. vite.config.ts), в production фронт и backend живут под одним origin.
const transport = createConnectTransport({ baseUrl: '/' });

export const companyClient = createClient(CompanyService, transport);
