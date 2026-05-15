import { useState, type FormEvent } from 'react';
import { AlertCircleIcon, InfoIcon, Loader2Icon } from 'lucide-react';
import { companyClient } from './api';
import type { GetCompanyResponse } from './gen/company/v1/company_pb';
import type { Stock } from './gen/company/v1/stock_pb';
import { formatValue, humanizeKey } from './format';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

type State =
  | { kind: 'idle' }
  | { kind: 'loading'; ticker: string }
  | { kind: 'ok'; ticker: string; data: GetCompanyResponse }
  | { kind: 'err'; ticker: string; error: string };

export function App() {
  const [ticker, setTicker] = useState('SBER');
  const [state, setState] = useState<State>({ kind: 'idle' });

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    const t = ticker.trim();
    if (!t) return;
    setState({ kind: 'loading', ticker: t });
    try {
      const data = await companyClient.getCompany({ ticker: t });
      setState({ kind: 'ok', ticker: t, data });
    } catch (err) {
      setState({
        kind: 'err',
        ticker: t,
        error: err instanceof Error ? err.message : String(err),
      });
    }
  }

  return (
    <div className="mx-auto max-w-7xl px-6 py-10">
      <header className="flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-muted-foreground">
            Financial Analyst
          </p>
          <h1 className="mt-1 text-3xl font-semibold tracking-tight">Карточка эмитента</h1>
        </div>
        <form onSubmit={onSubmit} className="flex items-center gap-2">
          <Input
            value={ticker}
            onChange={(e) => setTicker(e.target.value.toUpperCase())}
            placeholder="SBER"
            className="w-44 font-mono"
            spellCheck={false}
          />
          <Button type="submit" disabled={state.kind === 'loading'}>
            {state.kind === 'loading' && <Loader2Icon className="animate-spin" />}
            Загрузить
          </Button>
        </form>
      </header>

      <main className="mt-8">
        {state.kind === 'idle' && (
          <Alert>
            <InfoIcon />
            <AlertTitle>Начни с тикера</AlertTitle>
            <AlertDescription>
              Введите тикер (SBER, GAZP, LKOH, YDEX) и нажмите «Загрузить», чтобы увидеть все
              секции агрегата Company.
            </AlertDescription>
          </Alert>
        )}

        {state.kind === 'loading' && (
          <div className="space-y-3">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-64 w-full" />
          </div>
        )}

        {state.kind === 'err' && (
          <Alert variant="destructive">
            <AlertCircleIcon />
            <AlertTitle>Ошибка запроса</AlertTitle>
            <AlertDescription className="font-mono">{state.error}</AlertDescription>
          </Alert>
        )}

        {state.kind === 'ok' && <CompanyView data={state.data} />}
      </main>
    </div>
  );
}

type SectionDef = {
  value: string;
  label: string;
  description: string;
  kind: 'kv' | 'rows';
  data: Record<string, unknown> | Record<string, unknown>[] | undefined;
};

function CompanyView({ data }: { data: GetCompanyResponse }) {
  const company = data.company;
  if (!company) {
    return (
      <Alert variant="destructive">
        <AlertCircleIcon />
        <AlertTitle>Пустой ответ</AlertTitle>
        <AlertDescription>Сервер не вернул company.</AlertDescription>
      </Alert>
    );
  }
  const stock: Stock | undefined = company.stock;

  const sections: SectionDef[] = [
    {
      value: 'description',
      label: 'Description',
      description: 'Идентификаторы выпуска, классификация и режимы торгов',
      kind: 'kv',
      data: company.securityDescription,
    },
    {
      value: 'info',
      label: 'Info',
      description: 'Карточка эмитента (GICS, страна, описание)',
      kind: 'kv',
      data: stock?.info,
    },
    {
      value: 'summary',
      label: 'Summary',
      description: 'Сводные метрики одной строкой',
      kind: 'kv',
      data: stock?.summary,
    },
    {
      value: 'ratios',
      label: 'Ratios',
      description: 'Мультипликаторы по отчётным периодам',
      kind: 'rows',
      data: stock?.ratios,
    },
    {
      value: 'reports',
      label: 'Reports',
      description: 'Публикации финансовой отчётности',
      kind: 'rows',
      data: stock?.reports,
    },
    {
      value: 'dividends',
      label: 'Dividends',
      description: 'История и прогноз дивидендных выплат',
      kind: 'rows',
      data: stock?.dividends,
    },
    {
      value: 'ideas',
      label: 'Ideas',
      description: 'Инвест-идеи аналитиков',
      kind: 'rows',
      data: stock?.ideas,
    },
    {
      value: 'insiders',
      label: 'Insiders',
      description: 'Сделки инсайдеров',
      kind: 'rows',
      data: stock?.insiderTransactions,
    },
    {
      value: 'operations',
      label: 'Operations',
      description: 'Операционные метрики по периодам',
      kind: 'rows',
      data: stock?.operations,
    },
    {
      value: 'owners',
      label: 'Owners',
      description: 'Структура акционеров по датам среза',
      kind: 'rows',
      data: stock?.owners,
    },
    {
      value: 'shares',
      label: 'Shares',
      description: 'Количество выпущенных акций по датам среза',
      kind: 'rows',
      data: stock?.shares,
    },
  ];

  return (
    <Tabs defaultValue="description" className="gap-6">
      <TabsList className="flex h-auto w-full flex-wrap justify-start">
        {sections.map((s) => (
          <TabsTrigger key={s.value} value={s.value}>
            {s.label}
            {s.kind === 'rows' && Array.isArray(s.data) && s.data.length > 0 && (
              <Badge variant="secondary" className="ml-2">
                {s.data.length}
              </Badge>
            )}
          </TabsTrigger>
        ))}
      </TabsList>

      {sections.map((s) => (
        <TabsContent key={s.value} value={s.value}>
          <Card>
            <CardHeader className="border-b">
              <CardTitle>{s.label}</CardTitle>
              <CardDescription>{s.description}</CardDescription>
            </CardHeader>
            <CardContent>
              {s.kind === 'kv' ? (
                <KVTable data={s.data as Record<string, unknown> | undefined} />
              ) : (
                <RowsTable rows={s.data as Record<string, unknown>[] | undefined} />
              )}
            </CardContent>
          </Card>
        </TabsContent>
      ))}
    </Tabs>
  );
}

function KVTable({ data }: { data: Record<string, unknown> | undefined }) {
  if (!data || Object.keys(data).length === 0) return <Empty />;
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-1/3">Поле</TableHead>
          <TableHead>Значение</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {Object.entries(data).map(([key, value]) => (
          <TableRow key={key}>
            <TableCell className="font-medium text-muted-foreground">
              {humanizeKey(key)}
            </TableCell>
            <TableCell className="font-mono tabular-nums">{formatValue(value)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function RowsTable({ rows }: { rows: Record<string, unknown>[] | undefined }) {
  if (!rows || rows.length === 0) return <Empty />;
  const columns = Array.from(
    rows.reduce<Set<string>>((acc, row) => {
      Object.keys(row).forEach((k) => acc.add(k));
      return acc;
    }, new Set()),
  );
  return (
    <div className="-mx-4 overflow-x-auto px-4">
      <Table>
        <TableHeader>
          <TableRow>
            {columns.map((col) => (
              <TableHead key={col} className="whitespace-nowrap">
                {humanizeKey(col)}
              </TableHead>
            ))}
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row, i) => (
            <TableRow key={i}>
              {columns.map((col) => (
                <TableCell
                  key={col}
                  className="whitespace-nowrap font-mono text-xs tabular-nums"
                >
                  {formatValue(row[col])}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

function Empty() {
  return (
    <Alert>
      <InfoIcon />
      <AlertTitle>Пусто</AlertTitle>
      <AlertDescription>Источник не вернул данных для этой секции.</AlertDescription>
    </Alert>
  );
}
