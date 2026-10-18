import { useAppStore } from '@/store/useAppStore';

describe('useAppStore', () => {
  beforeEach(() => {
    useAppStore.setState({
      connection: null,
      apps: [],
      deployments: [],
      selectedApp: null,
      isLoading: false,
    });
  });

  test('sets connection correctly', () => {
    expect(useAppStore.getState().connection).toBeNull();

    useAppStore.getState().setConnection({
      url: 'https://demo.codedock.dev',
      token: 'cdk_secret_123',
      userEmail: 'dev@codedock.dev',
      userName: 'Developer',
      version: '1.0.0',
    });

    const conn = useAppStore.getState().connection;
    expect(conn?.url).toBe('https://demo.codedock.dev');
    expect(conn?.token).toBe('cdk_secret_123');
    expect(conn?.userEmail).toBe('dev@codedock.dev');
  });

  test('sets apps and selected app', () => {
    const mockApp = {
      id: 'app-1',
      name: 'api-service',
      status: 'running',
    };

    useAppStore.getState().setApps([mockApp]);
    expect(useAppStore.getState().apps).toHaveLength(1);

    useAppStore.getState().setSelectedApp(mockApp);
    expect(useAppStore.getState().selectedApp?.id).toBe('app-1');
  });
});
