import { infrastructureComparisons } from './infrastructure-comparisons';
import { managedComparisons } from './managed-comparisons';
export const comparisons = { ...infrastructureComparisons, ...managedComparisons };
